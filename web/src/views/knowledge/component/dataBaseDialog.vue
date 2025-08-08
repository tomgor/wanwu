<template>
  <el-dialog
    title="元数据管理"
    :visible.sync="dialogVisible"
    width="40%"
    :before-close="handleClose"
  >
    <div>
        <el-table
        :data="tableData"
        style="width: 100%">
        <el-table-column
            prop="key"
            label="key">
        </el-table-column>
        <el-table-column
            prop="value"
            label="value">
            <template slot-scope="scope">
                <span v-if="!scope.row.showInput">{{scope.row.value}}</span>
                <el-input v-model="scope.row.value" v-else @blur="handleBlur(scope.row)"></el-input>
            </template>
        </el-table-column>
        <el-table-column
            label="操作">
            <template slot-scope="scope">
                <el-button type="text" size="small" @click="editItem(scope.row)">编辑</el-button>
                <!-- <el-button type="text" size="small" @click="delItem(scope.$index)">删除</el-button> -->
            </template>
        </el-table-column>
        </el-table>
    </div>
    <span
      slot="footer"
      class="dialog-footer"
    >
      <el-button
        type="primary"
        @click="submitDialog"
      >确 定</el-button>
    </span>
  </el-dialog>
</template>
<script>
import { updateDocMeta} from "@/api/knowledge";
export default {
  data() {
    return {
      dialogVisible: false,
      tableData: [],
      docId: "",
    };
  },
  methods: {
    submitDialog() {
      this.tableData.forEach(i => delete i.showInput);
      const data = {
        MetaDataList:this.tableData,
        docId:this.docId
      }
      updateDocMeta(data).then(res =>{
        if(res.code === 0){
            this.$message.success('修改成功')
            this.$emit('updateData')
            this.dialogVisible = false;
        }
      })
      
    },
    editItem(n){
        n.showInput = true;
    },
    handleBlur(n){
        n.showInput = false;
    },
    delItem(index){
        this.tableData.splice(index,1);
    },
    showDiglog(data,id) {
      this.dialogVisible = true;
      this.docId = id;
      this.tableData = data.map(item => ({
          ...item,
          showInput: false
        }));
    },
    handleClose() {
      this.dialogVisible = false;
    }
  },
};
</script>
<style lang="scss" scoped>
/deep/ {
  .el-dialog__body {
    padding: 5px 20px !important;
  }
}
</style>