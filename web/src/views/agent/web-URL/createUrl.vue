<template>
    <div>
        <div>
            <el-button type="primary" icon="el-icon-plus" size="mini" @click="showDialog">创建</el-button>
            <el-table
                :data="tableData"
                style="width: 100%;margin-top:15px;"
                :header-cell-style="{ background: '#F9F9F9', color: '#999999' }"
              >
                <el-table-column
                  prop="name"
                  label="应用名称"
                >
                </el-table-column>
                <el-table-column
                  prop="suffix"
                  label="访问Url"
                >
                <template slot-scope="scope">
                    <span>{{scope.row.suffix}}</span>
                    <span class="el-icon-copy-document"></span>
                </template>
                </el-table-column>
                <el-table-column
                  prop="expiredAt"
                  label="过期时间"
                >
                </el-table-column>
                <el-table-column
                  prop="status"
                  label="状态"
                >
                <template slot-scope="scope">
                    <el-switch
                        v-model="scope.row.status"
                        active-color="#384BF7"
                        >
                    </el-switch>
                </template>
                </el-table-column>
                <el-table-column
                  :label="$t('knowledgeManage.operate')"
                  width="260"
                >
                  <template slot-scope="scope">
                    <el-button
                      size="mini"
                      round
                      @click="showDialog(scope.row)"
                    >编辑</el-button>
                    <el-button
                      size="mini"
                      round
                      @click="handleDel(scope.row)"
                    >{{$t('common.button.delete')}}</el-button>
                  </template>
                </el-table-column>
              </el-table>
        </div>
        <el-dialog
          :title="title"
          :visible.sync="dialogVisible"
          width="40%"
          :before-close="handleClose">
          <el-form ref="form" :model="form">
            <el-form-item label="应用名称" prop="name">
              <el-input v-model="form.name"></el-input>
            </el-form-item>
            <el-form-item label="过期时间">
              <el-input v-model="form.expiredAt"></el-input>
            </el-form-item>
            <!-- <el-form-item label="知识库出处详情">
              <el-switch
                v-model="value"
                active-color="#384BF7">
              </el-switch>
            </el-form-item> -->
            <!-- <el-form-item label="工作流详情">
              <el-switch
                v-model="value"
                active-color="#384BF7">
              </el-switch>
            </el-form-item> -->
            <el-form-item label="版权">
              <el-input v-model="form.copyright"></el-input>
              <el-switch
                v-model="form.copyrightEnable"
                active-color="#384BF7">
              </el-switch>
            </el-form-item>
            <el-form-item label="隐私协议">
              <el-input v-model="form.privacyPolicy"></el-input>
              <el-switch
                v-model="form.privacyPolicyEnable"
                active-color="#384BF7">
              </el-switch>
            </el-form-item>
            <el-form-item label="免责声明">
              <el-input v-model="form.disclaimer" type="textarea" :rows="2"></el-input>
              <el-switch
                v-model="form.disclaimerEnable"
                active-color="#384BF7">
              </el-switch>
            </el-form-item>
          </el-form>
          <span slot="footer" class="dialog-footer">
            <el-button @click="dialogVisible = false">取 消</el-button>
            <el-button type="primary" @click="dialogVisible = false">确 定</el-button>
          </span>
        </el-dialog>
    </div>
</template>
<script>
import { getOpenurl,delOpenurl } from "@/api/agent";
export default {
    props:['appId','appType'],
    data(){
        return{
          form:{
            appId:'',
            appType:'',
            copyright:'',
            copyrightEnable:false,
            disclaimer:'',
            disclaimerEnable:false,
            expiredAt:'',
            name:'',
            privacyPolicy:'',
            privacyPolicyEnable:false
          },
          title:'创建URL',
          dialogVisible:false,
          tableData:[]
        }
    },
    created(){
      this.getList();
    },
    methods:{
       getList(){
          getOpenurl({appId:this.appId,appType:this.appType}).then(res =>{
            if(res.code === 0){
              this.tableData = res.data.list || []
            }
          }).catch(() =>{

          })
       },
      showDialog(row=null){
        this.dialogVisible = true
        if(row === null){
          this.title = '创建URL';
        }else{
          this.title = '编辑URL';
        }
      },
      handleDel(row){
          this.$confirm(
          "确定要删除当前访问URL吗？",
          this.$t("knowledgeManage.tip"),
          {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning",
          }
        )
          .then(() => {
            delOpenurl({ urlId: row.urlId }).then((res) => {
              if (res.code === 0) {
                this.$message.success("删除成功");
                this.getList();
              }
            });
          })
          .catch((error) => {
            this.getList();
          });
        },
      handleClose(){
        this.dialogVisible = false;
      }
    }
}
</script>