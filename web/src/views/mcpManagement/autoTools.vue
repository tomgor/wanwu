<template>
  <div class="mcp-content-box customize">
    <div class="mcp-content">
      <div class="card-search card-search-cust">
        <div>
          <p class="card-search-des">
           创建自定义工具
          </p>
        </div>
        <div>
          <search-input placeholder="请输入名称搜索" ref="searchInput" @handleSearch="handleSearch" />
        </div>
      </div>

      <div class="card-box">
        <div class="card card-item-create">
          <div class="app-card-create" @click="handleAddMCP">
            <div class="create-img-wrap">
              <img class="create-type" src="@/assets/imgs/create_tools.svg" alt="" />
              <img class="create-img" src="@/assets/imgs/create_icon.png" alt="" />
              <div class="create-filter"></div>
            </div>
            <span>创建自定义工具</span>
          </div>
        </div>
        <div
          v-if="list && list.length"
          class="card"
          v-for="(item, index) in list"
          :key="index"
          @click.stop="handleClick(item)"
        >
          <div class="card-title">
            <img class="card-logo" src="@/assets/imgs/toolImg.svg" />
            <div class="mcp_detailBox">
              <span class="mcp_name">{{ item.name }}</span>
              <span class="mcp_from">
                <label>
                  {{ item.from }}
                </label>
              </span>
            </div>
            <i
              class="el-icon-delete-solid"
              @click.stop="handleDelete(item)"
            ></i>
          </div>
          <div class="card-des">{{ item.desc }}</div>
        </div>
      </div>
      <el-empty class="noData" v-if="!(list && list.length)" :description="$t('common.noData')"></el-empty>
    </div>
    <addDialog  ref="addDialog"></addDialog>
  </div>
</template>
<script>
import addDialog from "./addToolDialog.vue";
import SearchInput from "@/components/searchInput.vue"
import { getList, setDelete } from "@/api/mcp";
export default {
  components: { SearchInput, addDialog },
  data() {
    return {
      basePath: this.$basePath,
      is: true, // 是否第一次进来
      list: [],
    };
  },
  mounted() {
    this.is = true;
    this.fetchList((list) => {
      if (list.length) {
        this.is = false
      }
    })
  },
  methods: {
    radioChange(val) {
      this.handleSearch()
    },
    handleSearch() {
      if (this.is !== false) {
        if (this.list && this.list.length) {
          this.is = false
        } else {
          this.is = true
        }
      }
      this.fetchList()
    },
    fetchList(cb) {
      const searchInput = this.$refs.searchInput
      const params = {
        name: searchInput.value,
      }
      getList(params)
        .then((res) => {
          this.list = res.data.list || []
          cb && cb(this.list)
        })
        .catch(() => {})
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
    handleClick(val) {
      // smcpSquareId 有值 mcp广场, 否则自定义
      this.$router.push({path: `/mcp/detail/custom?mcpId=${val.mcpId}&mcpSquareId=${val.mcpSquareId}`})
    },
    handleAddMCP() {
      this.$refs.addDialog.showDialog();
    },
    handleDelete(item) {
      this.$confirm(
        "确定要删除 <span style='font-weight: bold;'>" + item.name + "</span> 该服务吗？",
        "提示",
        {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          dangerouslyUseHTMLString: true,
          type: "warning",
          center: true,
        }
      ).then(async () => {
        setDelete({
          mcpId: item.mcpId,
        }).then((res) => {
          if (res.code === 0) {
            this.$message.success("删除成功")
            this.init();
          } else {
            this.$message.error( res.msg || '删除失败')
          }
        })
      })
    },
  },
};
</script>
<style lang="scss">
.card-search-cust {
  text-align: left !important;

  .radio-box {
    margin: 20px 0 0 0 !important;
  }
}
.card-logo{
  width: 40px;
  height: 40px;
  object-fit: cover;
}
.mcp-content-box .noData {
  width: 100%;
  text-align: center;
  margin-top: -60px;
  /deep/ .el-empty__description p {
    color: #B3B1BC;
  }
}
</style>