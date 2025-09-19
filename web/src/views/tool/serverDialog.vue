<template>
  <div class="add-dialog">
    <el-dialog
        :title="title"
        :visible.sync="dialogVisible"
        width="50%"
        :show-close="false"
        :close-on-click-modal="false"
    >
      <div>
        <el-form
            :model="ruleForm"
            :rules="rules"
            ref="ruleForm"
            label-width="130px"
        >
          <el-form-item label="服务名称" prop="name">
            <el-input v-model="ruleForm.name"></el-input>
          </el-form-item>
          <el-form-item label="服务描述" prop="desc">
            <el-input v-model="ruleForm.desc"></el-input>
          </el-form-item>
          <el-form-item label="绑定应用" style="text-align: right">
            <el-button
                size="mini"
                @click="dialogBindVisible = true"
            >
              管理
            </el-button>
            <div class="api-list">
              <el-table
                  :data="bindList"
                  border
                  size="mini"
                  class="api-table"
                  :header-cell-style="{ textAlign: 'center' }"
              >
                <el-table-column
                    prop="methodName"
                    label="显示名称">
                  <template #default="{ row }">
                    <el-input
                        v-model="row.methodName"
                        size="mini"
                        @input="val => (row.methodName = val.replace(/[^a-zA-Z]/g, ''))"
                        @blur="handleBlur(row)"
                        :disabled="!row.editable"
                        placeholder="仅英语"/>
                  </template>
                </el-table-column>
                <el-table-column
                    prop="name"
                    label="应用"/>
                <el-table-column
                    prop="desc"
                    label="应用描述"/>
                <el-table-column
                    label="操作"
                    align="center"
                >
                  <template #default="{ row }">
                    <i class="el-icon-edit-outline table-opera-icon"
                       style="margin-right: 20px;"
                       @click="editItem(row)"/>
                    <i class="el-icon-delete table-opera-icon"
                       @click="delItem(row)"/>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-form-item>
        </el-form>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button @click="handleClose" size="mini">取 消</el-button>
        <el-button
            type="primary"
            size="mini"
            @click="submitForm('ruleForm')"
            :loading="publishLoading"
        >
          确定
        </el-button>
      </span>
      <el-dialog
          title="绑定应用"
          :visible.sync="dialogBindVisible"
          width="50%"
          :show-close="false"
          :close-on-click-modal="false"
          append-to-body
      >
        <div class="vertical-transfer">
          <el-transfer
              filterable
              filter-placeholder="搜索应用"
              :titles="['可选择', '已选择']"
              v-model="choices"
              :data="appList"
              :render-content="renderFunc"
              @change="handleChange">
          </el-transfer>
        </div>
        <div style="margin-top: 20px; text-align: right">
          <el-button
              size="mini"
              @click="addTool">
            +导入openapi
            <el-tooltip
                effect="dark"
                content="导入的openapi不会保存至应用。若想永久保存，方便后续调用，可在工具广场-自定义工具模块，将openapi添加为自定义工具。"
                placement="top-start"
            >
              <i class="el-icon-question" style="margin-left: 5px; cursor: pointer;"></i>
            </el-tooltip>
          </el-button>
        </div>
        <span slot="footer" class="dialog-footer">
          <el-button @click="handleBindCancel" size="mini">取 消</el-button>
          <el-button
              type="primary"
              size="mini"
              @click="submitBindForm"
          >
            确定
          </el-button>
        </span>
      </el-dialog>
    </el-dialog>
    <addToolDialog ref="addToolDialog" @addOpenapi="handleOpenapi"/>
  </div>
</template>
<script>
import addToolDialog from './addToolDialog'
import {editServer, addServer, getServerBind, getAppList} from "@/api/mcp.js";

const appTypeMap = {
  agent: '智能体',
  rag: '文本问答',
  workflow: '工作流',
  tool: '自定义工具',
  openapi: 'openapi'
};

export default {
  components: {
    addToolDialog
  },
  data() {
    return {
      dialogVisible: false,
      dialogBindVisible: false,
      title: '',
      appList: [],
      choices: [],
      choicesOrigin: [],
      bindList: [],
      openapiList: [],
      ruleForm: {
        MCPServerId: "",
        name: "",
        desc: "",
      },
      rules: {
        name: [{required: true, message: "请输入服务名称", trigger: "blur"}],
        desc: [
          {required: true, message: "请输入服务描述", trigger: "blur"},
        ],
      },
      publishLoading: false
    };
  },
  methods: {
    showDialog(mcpServerId) {
      this.dialogVisible = true
      getAppList().then((res) => {
        const appList = (res.data && res.data.list) || []
        this.appList = appList.map((item, index) => ({
          key: index + this.appList.length,
          label: item.name,
          ...item,
          flag: '0'
        }))
        if (mcpServerId) {
          this.ruleForm = {
            MCPServerId: mcpServerId,
            name: name,
            desc: desc
          }
          this.title = '修改自定义工具'
          const params = {
            mcpServerId: mcpServerId
          }
          getServerBind(params)
              .then((res) => {
                const {apps, desc, name} = res.data
                this.ruleForm = {
                  MCPServerId: mcpServerId,
                  name: name,
                  desc: desc
                }
                apps.forEach((app) => {
                  const foundIndex = this.appList.findIndex((item) => item.label === app.name)
                  if (foundIndex !== -1) {
                    this.choicesOrigin.push(this.appList[foundIndex].key)
                  }
                  this.bindList.push({
                    ...app,
                    editable: true,
                    flag: '0',
                    key: this.appList[foundIndex].key
                  })
                })
                this.choices = [...this.choicesOrigin]
              })
        } else this.title = '新增自定义工具'
      })
    },
    addTool() {
      this.$refs.addToolDialog.showDialog('', false, "导入openapi")
    },
    renderFunc(h, option) {
      return h('div', {
        style: {
          display: 'flex',
          alignItems: 'start'
        }
      }, [
        option.avatar && option.avatar.path ?
            h('img', {
              style: {
                width: '16px',
                height: '16px',
                marginTop: '8px',
                borderRadius: '50%'
              },
              attrs: {src: `${this.basePath}/user/api/${option.avatar.path}`},
            }) :
            h('img', {
              style: {
                width: '16px',
                height: '16px',
                marginTop: '8px',
                borderRadius: '50%'
              },
              attrs: {src: require('@/assets/imgs/toolImg.png')},
            }),
        h('span', {style: {marginLeft: '10px'}}, option.label),
        h('span', {
          style: {
            height: '24px',
            marginTop: '4px',
            marginLeft: '10px',
            lineHeight: '1.4',
            paddingLeft: '4px',
            paddingRight: '4px',
            backgroundColor: 'transparent',
            border: '1px solid #ccc',
            borderRadius: '4px',
            color: '#333'
          }
        }, appTypeMap[option.appType] || '未知类型')
      ]);
    },
    handleOpenapi(params) {
      this.openapiList.push(params)
      this.appList.push({
        key: this.appList.length,
        label: params.name,
        name: params.name,
        flag: '1',
        desc: params.description,
        appType: 'openapi',
        created: true
      })
      this.choices.push(this.appList.length - 1)
    },
    handleChange(value, direction, movedKeys) {
      if (direction === 'left') {
        this.appList = this.appList.filter((app) => {
          return !(movedKeys.includes(app.key) && app.flag === '1');
        });
      }
    },
    handleClose() {
      this.dialogVisible = false
      this.$emit("handleClose", false)
      this.$refs.ruleForm.resetFields()
      this.$refs.ruleForm.clearValidate()
      this.bindList = []
      this.choices = []
      this.choicesOrigin = []
      this.appList = []
      this.openapiList = []
      this.ruleForm = {
        MCPServerId: "",
        name: "",
        desc: "",
      }
    },
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.publishLoading = true
          const params = {
            ...this.ruleForm,
            apps: this.bindList.apps
          }
          if (this.ruleForm.MCPServerId) {
            editServer(params).then((res) => {
              if (res.code === 0) {
                this.$message.success("发布成功")
                this.publishLoading = false
                this.$emit("handleFetch", false)
                this.handleClose()
              }
            }).catch(() => this.publishLoading = false)
          } else {

          }
        }
      });
    },
    editItem(n) {
      n.editable = !n.editable;
      const index = this.bindList.indexOf(n);
      if (index !== -1) {
        this.appList.forEach((item) => {
          if (item.key === n.key) {
            item.methodName = n.methodName;
          }
        })
      }
    },
    handleBlur(n) {
      n.editable = false;
      const index = this.bindList.indexOf(n);
      if (index !== -1) {
        this.appList.forEach((item) => {
          if (item.key === n.key) {
            item.methodName = n.methodName;
          }
        })
      }
    },
    delItem(n) {
      const index = this.bindList.indexOf(n);
      if (index !== -1) {
        this.bindList.splice(index, 1);
        this.choices = this.choices.filter((key) => key !== n.key);
        this.appList = this.appList.filter((item) => !(item.key === n.key && item.flag === '1'));
      }
    },

    handleBindCancel() {
      this.dialogBindVisible = false;
      this.choices = [...this.choicesOrigin]
      this.appList = this.appList.filter((item) => item.created === true)
    },
    submitBindForm() {
      this.dialogBindVisible = false;
      this.choicesOrigin = [...this.choices]
      this.bindList = this.choices.map((key) => {
        const item = this.appList.find((app) => app.key === key);
        console.log(item)
        return {
          ...item,
          editable: true
        }
      })
    },
  },
};
</script>
<style lang="scss" scoped>
.add-dialog {
  .el-button.is-disabled {
    &:active {
      background: transparent !important;
      border-color: #ebeef5 !important;
    }
  }

  .mcpList {
    list-style: none;

    li {
      padding: 10px;
      border: 1px solid #ddd;
      border-radius: 5px;
      margin-bottom: 10px;
      background: #fff;
    }
  }
}

.table-opera-icon {
  font-size: 18px;
  cursor: pointer;
  color: #384BF7;
}

.api-list {
  .api-table /deep/ .el-table__body tr td,
  .api-table /deep/ .el-table__header tr th,
  .api-table /deep/ .el-table__body tr:hover > td {
    background-color: transparent !important;
  }
}

.vertical-transfer {
  /deep/ .el-transfer {
    display: flex;
    justify-content: space-between;
  }

  /deep/ .el-transfer-panel {
    width: 45%;
  }

  /deep/ .el-transfer__buttons {
    display: flex;
    flex-direction: column-reverse;
    justify-content: center;
    align-items: center;
    padding: 0;
  }

  /deep/ .el-transfer__buttons .el-button {
    margin: 20px 10px;
    width: 40px;
    padding: 8px 15px;
  }
}
</style>