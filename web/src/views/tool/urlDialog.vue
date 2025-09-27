<template>
  <div class="add-dialog">
    <el-dialog
        title="MCP服务地址"
        :visible.sync="dialogVisible"
        width="50%"
        @close="dialogVisible = false"
    >
      <div>
        <el-form
            :model="ruleForm"
            ref="ruleForm"
            label-width="130px"
        >
          <el-form-item label="服务地址（SSE）">
            <el-input v-model="ruleForm.sseUrl" disabled="disabled"/>
            <copyIcon :icon-class="ruleForm.sseUrl"/>
          </el-form-item>
          <el-form-item label="请求示例">
            <el-input
                type="textarea"
                rows="5"
                v-model="ruleForm.example"
                disabled="disabled"
            ></el-input>
            <copyIcon :icon-class="ruleForm.example"/>
          </el-form-item>
        </el-form>
      </div>
    </el-dialog>
  </div>
</template>
<script>
import {getServerUrl} from "@/api/mcp.js";
import CopyIcon from "@/components/copyIcon.vue";

export default {
  components: {CopyIcon},
  data() {
    return {
      dialogVisible: false,
      ruleForm: {
        example: "",
        sseUrl: ""
      },
    };
  },
  methods: {
    showDialog(mcpServerId) {
      this.dialogVisible = true
      const params = {
        mcpServerId: mcpServerId,
      }
      getServerUrl(params).then((res) => {
        this.ruleForm = res.data
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
</style>