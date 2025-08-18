<template>
    <div>
        <div>
            <el-button type="primary" icon="el-icon-plus" size="mini">创建</el-button>
            <el-table
                :data="tableData"
                style="width: 100%;margin-top:15px;"
                :header-cell-style="{ background: '#F9F9F9', color: '#999999' }"
              >
                <el-table-column
                  prop="docType"
                  label="应用名称"
                >
                </el-table-column>
                <el-table-column
                  prop="uploadTime"
                  label="访问Url"
                >
                <template slot-scope="scope">
                    <span></span>
                    <span class="el-icon-copy-document"></span>
                </template>
                </el-table-column>
                <el-table-column
                  prop="docType"
                  label="过期时间"
                >
                </el-table-column>
                <el-table-column
                  prop="docType"
                  label="状态"
                >
                <template slot-scope="scope">
                    <el-switch
                        v-model="value"
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
                      @click="handleEdit(scope.row)"
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
    </div>
</template>
<script>
export default {
    data(){
        return{
            tableData:[]
        }
    },
    methods:{
        handleEdit(row){

        },
        handleDel(row){
            
        }
    }
}
</script>