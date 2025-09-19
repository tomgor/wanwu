<template>
  <div class="add-dialog">
    <el-dialog
        title="MCP服务地址"
        :visible.sync="dialogVisible"
        width="50%"
        :show-close="false"
        :close-on-click-modal="false"
    >
      <div>
        <el-form
            :model="ruleForm"
            ref="ruleForm"
            label-width="130px"
        >
          <el-form-item label="服务地址（SSE）" prop="name">
            <el-input v-model="ruleForm.name" disabled="disabled"/>
            <copyIcon :icon-class="ruleForm.name"/>
          </el-form-item>
          <el-form-item label="请求示例" prop="desc">
            <el-input
                type="textarea"
                rows="5"
                v-model="ruleForm.desc"
                disabled="disabled"
            ></el-input>
            <copyIcon :icon-class="ruleForm.desc"/>
          </el-form-item>
        </el-form>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button @click="handleCancel" size="mini">取 消</el-button>
        <el-button
            type="primary"
            size="mini"
            :disabled="mcpList.length === 0"
            @click="submitForm('ruleForm')"
            :loading="publishLoading"
        >
          确定
        </el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
import {} from "@/api/mcp.js";
import CopyIcon from "@/components/copyIcon.vue";

export default {
  components: {CopyIcon},
  props: ["dialogVisible"],
  data() {
    return {
      appList: [],
      bindList: [],
      mcpList: [],
      ruleForm: {
        name: "",
        from: "",
      },
      publishLoading: false
    };
  },
  methods: {
    fetchData() {

    },
    handleCancel() {
      this.$emit("handleClose", false);
    },
    submitForm(formName) {

    }
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
</style>