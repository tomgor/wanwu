<template>
  <div>
    <div class="header">
      <div class="header-api">
        <el-tag
          effect="plain"
          class="root-url"
        >API根地址</el-tag>
        {{apiURL}}
      </div>
      <el-button
        size="small"
        @click="openApiDialog"
        class="apikeyBtn"
      >
        <img :src="require('@/assets/imgs/apikey.png')" />
        API秘钥
      </el-button>
    </div>
    <el-table
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column
        label="密钥"
        prop="apiKey"
        width="300"
      />
      <el-table-column
        label="创建时间"
        prop="createdAt"
      />
      <el-table-column
        label="操作"
        width="200"
      >
        <template slot-scope="scope">
          <el-button
            size="mini"
            @click="handleCopy(scope.row) && copycb()"
          >复制</el-button>
          <el-button
            size="mini"
            @click="handleDelete(scope.row)"
          >删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <!-- apikey -->
    <ApiKeyDialog
      ref="apiKeyDialog"
      :appId="appId"
      :appType="'agent'"
    />
  </div>
</template>
<script>
import ApiKeyDialog from "../components/ApiKeyDialog";
import { createApiKey, delApiKey, getApiKeyList } from "@/api/appspace";
export default {
  components:{ApiKeyDialog},
  props: {
    appType: {
      type: String,
      required: true,
    },
    appId: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      apiURL: "",
      tableData: [],
    };
  },
  created() {
    this.getTableData();
  },
  methods: {
    openApiDialog(){
      this.$refs.apiKeyDialog.showDialog();
    },
    handleCopy(row) {
      let text = row.apiKey;
      var textareaEl = document.createElement("textarea");
      textareaEl.setAttribute("readonly", "readonly");
      textareaEl.value = text;
      document.body.appendChild(textareaEl);
      textareaEl.select();
      var res = document.execCommand("copy");
      document.body.removeChild(textareaEl);
      return res;
    },
    copycb() {
      this.$message.success("内容已复制到粘贴板");
    },
    handleCreate() {
      const data = { appId: this.appId, appType: this.appType };
      createApiKey(data).then((res) => {
        if (res.code === 0) {
          this.tableData.push(res.data);
        }
      });
    },
    getTableData() {
      const data = { appId: this.appId, appType: this.appType };
      getApiKeyList(data).then((res) => {
        if (res.code === 0) {
          this.tableData = res.data || [];
        }
      });
    },
    handleDelete(row) {
      this.$confirm(
        "确定要删除当前APIkey吗？",
        this.$t("knowledgeManage.tip"),
        {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        }
      )
        .then(() => {
          delApiKey({ apiId: row.apiId }).then((res) => {
            if (res.code === 0) {
              this.$message.success("删除成功");
              this.getTableData();
            }
          });
        })
        .catch((error) => {
          this.getTableData();
        });
    },
  },
};
</script>
<style lang="scss" scoped>
  .header{
    width:100%;
    display:flex;
    justify-content:flex-start;
    align-items:flex-start;
    height:60px;
    .header-api {
      padding: 6px 10px;
      box-shadow: 1px 2px 2px #ddd;
      background-color: #fff;
      border-radius: 6px;
      width:20%;
      .root-url {
        background-color: #eceefe;
        color: #384bf7;
        border: none;
      }
    }
    .apikeyBtn{
      margin-left:10px;
      border:1px solid #384bf7;
      padding:12px;
      display:flex;
      align-items:center;
    }
  }

</style>