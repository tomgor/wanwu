<template>
  <div class="mcp-detail">
    <el-drawer
      title="服务详情"
      :visible.sync="drawer"
      :direction="direction"
      :before-close="handleClose"
      :show-close="true"
      :size="'400px'"
    >
      <div class="detail_content" v-loading="loading">
        <h2>
          {{ detail.name }} <span>{{ detail.serverFrom }}</span>
          <!-- 【服务来源】 -->
        </h2>
        <div class="detail-des">
          {{ detail.description }}
        </div>
        <el-divider></el-divider>
        <h3>MCP Server URL</h3>
        <div class="detail-url">
          <el-input
            placeholder="请输入内容"
            v-model="detail.serverUrl"
            :disabled="true"
          >
            <el-button
              slot="append"
              icon="el-icon-copy-document"
              @click="handleCopy(detail.serverUrl)"
            ></el-button>
          </el-input>
        </div>
        <h3>工具列表</h3>
        <ul class="detail-mcpList">
          <li v-for="(item, index) in mcpList" :key="index">{{ item.name }}</li>
        </ul>
        <h3>其他 MCP 服务查看</h3>
        <ul class="detail-other" v-if="list.length > 0">
          <li
            v-for="(item, index) in list"
            :key="index"
            @click="handleClick(item)"
          >
            <div class="card">
              <div class="card-title">
                <label>{{ item.name }}</label>

                <span>{{ item.serverFrom }}</span>
              </div>
              <div class="card-des">{{ item.description }}</div>
            </div>
          </li>
        </ul>
      </div>
    </el-drawer>
  </div>
</template>
<script>
import { getDetail, getTools, getList } from "@/api/mcp";

export default {
  props: ["drawer", "id"],
  data() {
    return {
      loading: false,
      direction: "rtl",
      detail: {
        name: "",
        serverFrom: "",
        description: "",
      },
      mcpList: [], // 工具列表
      list: [], // 随机mcp服务列表
    };
  },
  created() {
    this.init(this.id);
    this.getMcpList();
  },
  methods: {
    init(id) {
      this.loading = true;
      getDetail({
        mcpId: id,
      })
        .then((res) => {
          this.detail = res.data;

          getTools({
            serverUrl: res.data.serverUrl,
          })
            .then((res) => {
              this.mcpList = res.data.tools;
              this.loading = false;
            })
            .catch((err) => {
              this.loading = false;
            });
        })
        .catch((err) => {
          this.loading = false;
        });
    },
    handleClose() {
      this.$emit("handleDetailClose", false);
    },
    getMcpList() {
      this.list = [];
      getList({
        pageNo: 1,
        pageSize: 999,
      })
        .then((res) => {
          let list = res.data.list;
          if (list.length > 5) {
            function getRandomInRange(min, max) {
              return Math.floor(Math.random() * (max - min + 1)) + min;
            }

            // 假设我们想在1到100之间生成5个随机整数
            const numbers = [];
            const min = 0;
            const max = list.length - 1;

            while (this.list.length < 5) {
              const randomNumber = getRandomInRange(min, max);
              if (!numbers.includes(randomNumber)) {
                // 确保不重复
                numbers.push(randomNumber);
                this.list.push(list[randomNumber]);
              }
            }
          } else {
            this.list = list;
          }
          console.log(this.list);
        })
        .catch((err) => {});
    },
    handleClick(item) {
      this.init(item.mcpId);
      this.getMcpList();
    },
    handleCopy(url) {
      var textareaEl = document.createElement("textarea");
      textareaEl.setAttribute("readonly", "readonly"); // 防止手机上弹出软键盘
      textareaEl.value = url;
      document.body.appendChild(textareaEl);
      textareaEl.select();
      var res = document.execCommand("copy");
      document.body.removeChild(textareaEl);
      this.$message({
        message: "复制成功",
        type: "success",
      });
      return res;
    },
  },
};
</script>
<style lang="scss" scoped>
.mcp-detail {
  .el-drawer__header > :first-child {
    font-size: 16px;
  }
  .el-drawer__header {
    margin-bottom: 10px;
  }
  .detail_content {
    padding: 20px;
    h2 {
      font-size: 16px;
      word-wrap: break-word;
      span {
        margin-left: 10px;
        font-size: 10px;
        background: $color;
        color: #fff;
        padding: 3px 5px;
        border-radius: 3px;
      }
    }
    .el-loading-mask {
      height: calc(100vh - 53px);
      background-color: rgba(255, 255, 255, 0.8);
    }
    .detail-des {
      margin-top: 10px;
      color: #585a73;
      font-size: 13px;
    }
    .detail-mcpList {
      max-height: 450px;
      overflow-y: auto;
      margin-bottom: 10px;
      li {
        padding: 10px;
        border: 1px solid #ddd;
        border-radius: 5px;
        margin-bottom: 10px;
      }
    }
    .detail-url {
      margin-bottom: 10px;
      .el-input {
        width: 100%;
      }
    }
    ul {
      list-style: none;
    }
    h3 {
      margin-bottom: 5px;
    }
    .card {
      position: relative;
      padding: 20px 16px;
      /*border: 1px solid #e8e9eb;*/
      border-radius: 12px;
      background: #fff;
      //box-shadow: 0 2px 2px #0000000a;
      display: flex;
      flex-direction: column;
      align-items: center;
      margin-bottom: 10px;
      //border: 1px solid $border_color;
      box-shadow: 0 1px 4px 0 rgba(0, 0, 0, 0.15);
      &:hover {
        cursor: pointer;
        border: 1px solid $border_color;
        box-shadow: 0 2px 8px #171a220d, 0 4px 16px #0000000f;
      }

      .card-title {
        display: flex;
        align-items: center;
        margin-bottom: 20px;
        width: 100%;

        label {
          color: #5d5d5d;
          font-weight: 700;
          font-size: 15px;
          max-width: 60%;
          display: inline-block;
          overflow: hidden;
          white-space: nowrap;
          text-overflow: ellipsis;
          margin-right: 2%;
        }
        span {
          max-width: 37%;
          display: inline-block;
          overflow: hidden;
          white-space: nowrap;
          text-overflow: ellipsis;
          padding: 0.22em 0.5em;
          font-size: 10px;
          color: #84868c;
          background: #f2f5f9;
          border-radius: 3px;
        }
      }
      .card-des {
        display: -webkit-box;
        text-overflow: ellipsis;
        overflow: hidden;
        -webkit-line-clamp: 2;
        line-clamp: 2;
        -webkit-box-orient: vertical;
        font-size: 12px;
        word-wrap: break-word;
        width: 100%;
        color: #999;
      }
    }
  }
}
</style>