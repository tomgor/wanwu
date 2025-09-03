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
          <div class="app-card-create" @click="handleAddMCP('')">
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
          @click.stop="handleClick(item.customToolId)"
        >
          <div class="card-title">
            <img class="card-logo" src="@/assets/imgs/toolImg.png" />
            <div class="mcp_detailBox">
              <span class="mcp_name">{{ item.name }}</span>
            </div>
            <el-dropdown
                placement="bottom">
              <span class="el-dropdown-link">
                <i class="el-icon-more"
                    @click.stop/>
              </span>
              <el-dropdown-menu slot="dropdown"  style="margin-top: -10px">
                <el-dropdown-item
                    @click.native="handleAddMCP(item.customToolId)">
                  编辑
                </el-dropdown-item>
                <el-dropdown-item
                    @click.native="handleDelete(item)">
                  删除
                </el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
          </div>
          <div class="card-des">{{ item.description }}</div>
        </div>
      </div>
      <el-empty class="noData" v-if="!(list && list.length)" :description="$t('common.noData')"></el-empty>
    </div>
    <addDialog  ref="addDialog" @handleClose="handleClose"></addDialog>
  </div>
</template>
<script>
import addDialog from "./addToolDialog.vue";
import SearchInput from "@/components/searchInput.vue"
import { getCustomList, deleteCustom } from "@/api/mcp";
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
      getCustomList(params)
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
    handleClick(customToolId) {
      this.$refs.addDialog.showDialog(customToolId, true);
    },
    handleAddMCP(customToolId) {
      this.$refs.addDialog.showDialog(customToolId, false);
    },
    handleDelete(item) {
      this.$confirm(
        "删除后，历史引用了本自定义工具的智能体将自动取消引用，且此操作不可撤回,确认删除吗？",
        "提示",
        {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          dangerouslyUseHTMLString: true,
          type: "warning",
          center: true,
        }
      ).then(async () => {
        deleteCustom({
          customToolId: item.customToolId,
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
    handleClose() {
      this.init()
    }
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