<template>
  <el-dialog
    title="元数据管理"
    :visible.sync="dialogVisible"
    width="40%"
    :before-close="handleClose"
  >
    <el-button
        @click="addItem"
    >
      + 创建元数据
    </el-button>
    <div>
        <el-table
        :data="filteredTableData"
        style="width: 100%;margin-top: 20px">
          <el-table-column
              prop="key"
              label="Key"
              align="center">
            <template #default="{ row }">
              <el-input
                  v-model="row.key"
                  placeholder="只能包含小写字母、数字和下划线，并且必须以小写字母开头"
                  clearable
                  :disabled="!row.editable || !row.created"
              />
            </template>
          </el-table-column>
          <el-table-column
              prop="dataType"
              label="类型"
              align="center">
            <template #default="{ row }">
              <el-select
                  v-model="row.dataType"
                  placeholder="请选择"
                  clearable
                  :disabled="!row.editable"
              >
                <el-option value="string" label="String"></el-option>
                <el-option value="number" label="Number"></el-option>
                <el-option value="time" label="Time"></el-option>
              </el-select>
            </template>
          </el-table-column>
          <el-table-column
              prop="value"
              label="Value"
              align="center">
            <template #default="{ row }">
              <el-input
                  v-model="row.value"
                  @blur="handleBlur(row)"
                  clearable
                  :disabled="!row.editable"
              />
            </template>
          </el-table-column>
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
    <span
      slot="footer"
      class="dialog-footer"
    >
      <el-button @click="dialogVisible = false">取 消</el-button>
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
  computed: {
    filteredTableData() {
      return this.tableData.filter(item => item.option !== "delete");
    }
  },
  data() {
    return {
      dialogVisible: false,
      tableData: [],
      docId: "",
    };
  },
  methods: {
    submitDialog() {
      this.tableData.forEach(i => {
        delete i.editable
        delete i.created
      });
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
    addItem(){
      this.tableData.push({
        key: '',
        dataType: '',
        value: [],
        option: 'add',
        editable: true,
        created: true,
      });
    },
    editItem(n){
        n.editable = !n.editable;
    },
    handleBlur(n){
        n.editable = false;
    },
    delItem(index){
      index.option = "delete";
    },
    showDiglog(data,id) {
      this.dialogVisible = true;
      this.docId = id;
      this.tableData = data.map(item => ({
          ...item,
          option: 'update',
          editable: false,
          created: false,
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
.table-opera-icon{
  font-size: 18px;
}
</style>