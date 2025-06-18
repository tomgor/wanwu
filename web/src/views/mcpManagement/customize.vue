<template>
  <div class="mcp-content-box customize">
    <div class="mcp-content">
      <div class="card-search card-search-cust">
        <div>
          <p class="card-search-des">
            集合各MCP server，即选即用。可支持在插件编排工具中使用。如使用自定义MCP，请在自定义TAB中进行添加和查看&nbsp;
          </p>
          <!--<div class="radio-box" >
            <el-radio-group v-model="source" @input="radioChange" size="mini">
              <el-radio label="">全部</el-radio>
              <el-radio label="1">来自我添加</el-radio>
              <el-radio label="2">来自广场发送</el-radio>
            </el-radio-group>
          </div>-->
        </div>
        <div>
          <search-input placeholder="请输入内容" ref="searchInput" @handleSearch="handleSearch" />
          <el-button size="mini" type="primary" @click="handleAddMCP">导入</el-button>
        </div>
      </div>

      <div class="card-box" v-if="list.length > 0">
        <div
          class="card"
          v-for="(item, index) in list"
          :key="index"
          @click.stop="handleClick(item)"
        >
          <div class="card-title">
            <svg-icon icon-class="mcp_server" class="editable--send" />
            <div class="mcp_detailBox">
              <span class="mcp_name">{{ item.name }}</span>
              <span class="mcp_from">
                <label>
                  {{ item.serverFrom }}
                </label>
              </span>
            </div>
            <i
                class="el-icon-delete-solid"
                @click.stop="handleDelete(item)"
            ></i>
          </div>
          <div class="card-des">{{ item.description }}</div>
        </div>
      </div>

      <div class="no-list" v-if="list.length === 0 && is">
        <div>
          <i class="el-icon-circle-plus-outline" @click="handleAddMCP"></i>
          <span>添加你的第一个MCP Server</span>
        </div>
      </div>
      <!-- </el-tab-pane>
      </el-tabs> -->
    </div>
    <addDialog :dialogVisible="addOpen" @handleClose="handleClose"></addDialog>
    <detail
      v-if="drawer"
      :drawer="drawer"
      :id="mcpid"
      @handleDetailClose="handleDetailClose"
    ></detail>
  </div>
</template>
<script>
import addDialog from "./addDialog.vue";
import detail from "./detailDrawer.vue";
import { getList, setDelete } from "@/api/mcp";
import SearchInput from "@/components/searchInput.vue"
export default {
  components: { SearchInput, addDialog, detail },
  data() {
    return {
      is: true, // 是否第一次进来
      mcpid: "",
      name: "",
      activeName: "cust",
      addOpen: false, // 自定义添加mcp开关
      drawer: false, // 详情开关
      list: [
        {
          "mcpId": "aa513b0f-b4b5-4db9-8502-8b7689afb6f4",
          "serverFrom": "测试",
          "serverUrl": "https://mcp.amap.com/sse?key=6d889bd6aa34bdd63a39c1197a00e377",
          "description": "测试",
          "name": "测试",
          "source": 1,
          "mcpSquareId": ""
        }
      ],
      source:''
    };
  },
  mounted() {
    this.is = true;
    this.fetchList((list) => {
      if (list.length > 0) {
        this.is = false;
      }
    })
  },
  methods: {
    radioChange(val) {
      this.handleSearch()
    },
    handleSearch() {
      if (this.is !== false) {
        if (this.list.length > 0) {
          this.is = false;
        } else {
          this.is = true;
        }
      }
      this.fetchList()
    },
    fetchList(cb) {
      const searchInput = this.$refs.searchInput
      let params = {
        name: searchInput.value,
        pageNo: 1,
        pageSize: 9999,
      }
      if (this.source !== '') {
        params.source = this.source
      }
      getList(params)
        .then((res) => {
          this.list = res.data.list;
          cb && cb(this.list)
        })
        .catch((err) => {});
    },
    init() {
      this.fetchList((list) => {
        if (!this.is) return
        if (list.length > 0) {
          this.is = false
        } else {
          this.is = true
        }
      })
    },
    handleClose(val) {
      this.addOpen = val;
      this.init();
    },
    /*handleClick(val) {
        this.drawer = true;
        this.mcpid = val.mcpId;
    },*/
    handleClick(val) {
      // source 1自定义 2mcp广场
      this.$router.push({path: `/mcp/detail/${val.mcpId}/${val.source}/${val.mcpSquareId || '0'}`})
    },
    handleAddMCP() {
      this.addOpen = true;
    },
    handleDetailClose(val) {
      this.drawer = val;
    },
    handleDelete(item) {
      this.$confirm(
        "确定要删除 <span style='font-weight: bold;'>" +
        item.name +
        "</span> 该服务吗？",
        "提示",
        {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          dangerouslyUseHTMLString: true,
          type: "warning",
          center: true,
        }
      )
        .then(async () => {
          setDelete({
            mcpId: item.mcpId,
          })
            .then((res) => {
              if (res.code == 0) {
                this.$message({
                  type: "success",
                  message: "删除成功",
                });
                this.init();
              } else {
                this.$message({
                  type: "error",
                  message: res.msg,
                });
              }
            })
            .catch((err) => {})
        })
        .catch(() => {})
    },
  },
};
</script>
<style lang="scss" scoped>
.card-search-cust {
  text-align: left !important;

  .radio-box {
    margin: 20px 0 0 0 !important;
  }
}

.customize {
  .card-box {
   /* max-height: calc(100vh - 130px);
    overflow-y: auto;*/
  }
}
</style>