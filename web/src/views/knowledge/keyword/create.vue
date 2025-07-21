<template>
    <div>
       <el-dialog
        :title="title"
        :visible.sync="dialogVisible"
        width="40%"
        :before-close="handleClose">
        <el-form ref="form" :model="form" label-width="120px">
            <el-form-item label="问题中的关键词">
                <el-input v-model="form.name"></el-input>
            </el-form-item>
            <el-form-item label="文档中的词语">
                <el-input v-model="form.name"></el-input>
            </el-form-item>
            <el-form-item label="选择知识库">
            <el-select
              v-model="form.knowledgeIdList"
              placeholder="请选择"
              multiple
              clearable
              filterable 
              style="width:100%;"
              @visible-change="visibleChange($event)"
            >
              <el-option
                v-for="item in knowledgeOptions"
                :key="item.knowledgeId"
                :label="item.name"
                :value="item.knowledgeId"
              >
              </el-option>
            </el-select>
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
import { getKnowledgeList} from "@/api/knowledge";
export default {
    data(){
        return{
            form:{},
            knowledgeOptions:[],
            title:'创建关键词',
            dialogVisible:false
        }
    },
    created() {
        this.getKnowledgeList();
    },
    methods:{
      async getKnowledgeList() {
            //获取文档知识分类
            const res = await getKnowledgeList({});
            if (res.code === 0) {
                this.knowledgeOptions = res.data.knowledgeList || [];
            } else {
                this.$message.error(res.message);
            }
        },
        visibleChange(val) {
            if (val) {
                this.getKnowledgeList();
            }
        },
        showDialog(){
            this.dialogVisible = true
        },
        handleClose(){
            this.dialogVisible = false
        }
    }
}
</script>