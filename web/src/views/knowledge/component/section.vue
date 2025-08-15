<template>
  <div class="section" v-loading="loading.itemStatus">
    <div class="title">
      <i
        class="el-icon-arrow-left"
        @click="$router.go(-1)"
        style="margin-right: 20px; font-size: 20px; cursor: pointer"
      ></i
      >{{ obj.name }}
    </div>
    <div class="container">
      <el-descriptions
        class="margin-top"
        title=""
        :column="3"
        :size="''"
        border
      >
        <el-descriptions-item :label="$t('knowledgeManage.fileName')">{{
          res.fileName
        }}</el-descriptions-item>
        <el-descriptions-item :label="$t('knowledgeManage.splitNum')">
          {{ res.segmentTotalNum }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('knowledgeManage.importTime')">{{
          res.uploadTime
        }}</el-descriptions-item>
        <el-descriptions-item :label="$t('knowledgeManage.chunkType')">{{
          Number(res.segmentType) === 0 ? $t('knowledgeManage.autoChunk') : $t('knowledgeManage.autoConfigChunk')
        }}</el-descriptions-item>
        <el-descriptions-item :label="$t('knowledgeManage.setMaxLength')">{{
          String(res.maxSegmentSize)
        }}</el-descriptions-item>
        <el-descriptions-item :label="$t('knowledgeManage.markSplit')">{{
          res.splitter
        }}</el-descriptions-item>
        <el-descriptions-item label="元数据">
          <template v-if="metaDataList && metaDataList.length > 0">
            <span
                v-for="(item, index) in metaDataList.slice(0, 3)"
                :key="index"
                class="metaItem"
            >
              {{ item.key }}: {{ item.dataType === 'time' ? formatTimestamp(item.value) : item.value }}
              <span v-if="index < metaDataList.slice(0, 3).length - 1">, </span>
            </span>
            <el-tooltip v-if="metaDataList.length > 3" :content="filterData(metaDataList.slice(3))" placement="bottom">
              <span class="metaItem">...</span>
            </el-tooltip>
          </template>
          <span v-else>无数据</span>
          <span class="el-icon-edit-outline editIcon" @click="showDatabase(metaDataList || [])" v-if="metaDataList"></span>
        </el-descriptions-item>
        <el-descriptions-item label="元数据规则">
          <template v-if="metaRuleList && metaRuleList.length > 0">
            <span v-for="(item, index) in metaRuleList.slice(0, 3)" :key="index" class="metaItem">
              {{ item.key }}: {{ item.rule }}<span v-if="index < metaRuleList.slice(0, 3).length - 1"> </span>
            </span>
            <el-tooltip v-if="metaRuleList.length > 3" :content="filterRule(metaRuleList.slice(3))" placement="bottom">
              <span class="metaItem">...</span>
            </el-tooltip>
          </template>
          <span v-else>无数据</span>
        </el-descriptions-item>
      </el-descriptions>

      <div class="btn">
        <el-button
          type="primary"
          @click="handleStatus('start')"
          size="mini"
          :loading="loading.start"
          >{{$t('knowledgeManage.allRun')}}</el-button
        >
        <el-button
          type="primary"
          @click="handleStatus('stop')"
          size="mini"
          :loading="loading.stop"
          >{{$t('knowledgeManage.allStop')}}</el-button
        >
      </div>

      <div class="card">
        <el-row :gutter="20" v-if="res.contentList.length > 0">
          <el-col
            :span="6"
            v-for="(item, index) in res.contentList"
            :key="index"
            class="card-box"
          >
            <el-card class="box-card">
              <div slot="header" class="clearfix">
                <span
                  >{{ $t('knowledgeManage.split')+":" + item.contentNum }}&nbsp;&nbsp;
                  <span style="font-size: 12px"
                    >{{$t('knowledgeManage.length')}}:{{ item.content.length }}{{$t('knowledgeManage.character')}}</span
                  >
                </span>

                <el-switch
                  style="float: right; padding: 3px 0"
                  v-model="item.available"
                  active-color="#384bf7"
                  @change="handleStatusChange(item, index)"
                >
                </el-switch>
              </div>
              <div class="text item" @click="handleClick(item, index)">
                {{ item.content }}
              </div>
            </el-card>
          </el-col>
        </el-row>
        <el-empty v-else :description="$t('knowledgeManage.noData')"></el-empty>
      </div>

      <div class="list-common" style="text-align: right">
        <el-pagination
          background
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="page.pageNo"
          :page-sizes="page.pageSizeList"
          :page-size="page.pageSize"
          layout="total, prev, pager, next, jumper"
          :total="page.total"
        >
        </el-pagination>
      </div>
    </div>

    <!-- 详情弹框 -->
    <el-dialog
      v-if="dialogVisible"
      :title="$t('knowledgeManage.detailView')"
      :visible.sync="dialogVisible"
      width="60%"
      :show-close="false"
      v-loading="loading.dialog"
    >
      <div slot="title">
        <span style="font-size: 16px">{{$t('knowledgeManage.detailView')}}</span>
        <el-switch
          @change="handleDetailStatusChange"
          style="float: right; padding: 3px 0"
          v-model="cardObj[0].available"
          active-color="#384bf7"
        >
        </el-switch>
      </div>
      <div>
        <el-table
          :data="cardObj"
          border
          style="width: 100%"
          :header-cell-style="{
            background: '#F9F9F9',
            color: '#999999',
          }"
        >
          <el-table-column
            prop="content"
            align="center"
            :render-header="renderHeader"
          >
          </el-table-column>
        </el-table>
      </div>

      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="handleClose">{{$t('knowledgeManage.close')}}</el-button>
      </span>
    </el-dialog>
    <dataBaseDialog ref="dataBase" @updateData="updateData"/>
  </div>
</template>
<script>
import { getSectionList,setSectionStatus } from "@/api/knowledge";
import dataBaseDialog from './dataBaseDialog'
export default {
  components:{dataBaseDialog},
  data() {
    return {
      dialogVisible: false,
      obj: {}, // 路由参数对象
      cardObj: [
        {
          available: false,
          content: "",
          contentId: "",
          len: 20,
        },
      ], // 单独卡片存储对象
      value: true,
      activeStatus: false,
      page: {
        pageNo: 1,
        pageSize: 8,
        pageSizeList: [10, 15, 20, 50],
        total: 0,
      },
      loading: {
        start: false,
        stop: false,
        itemStatus: false,
        dialog: false,
      },
      res: {
        contentList: [],
      },
      metaDataList: [],
      metaRuleList: []
    };
  },
  created() {
    this.obj = this.$route.query;
    this.getList();
  },
  methods: {
    updateData(){
      this.getList();
    },
    showDatabase(data){
      this.$refs.dataBase.showDiglog(data,this.obj.id)
    },
    filterData(data){
      const formattedString = data.map(item => {
        let value = item.value;
        // 如果是时间类型且值为时间戳，转换为日期字符串
        if (item.dataType === 'time' && typeof value === 'number') {
          value = this.formatTimestamp(value);
        }
        return `${item.key}:${value}`;
      }).join(", ");
      return formattedString;
    },
    formatTimestamp(timestamp) {
      const date = new Date(Number(timestamp));
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      const seconds = String(date.getSeconds()).padStart(2, '0');
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    },
    filterRule(rule){
      const formattedString = rule.map(item => `${item.key}:${item.rule}`).join(", ");
      return formattedString
    },
    getList() {
      this.loading.itemStatus = true;
      getSectionList({
        docId: this.obj.id,
        pageNo: this.page.pageNo,
        pageSize:this.page.pageSize
      })
        .then((res) => {
          this.loading.itemStatus = false;
          this.res = res.data;
          this.page.total = this.res.segmentTotalNum;
          this.metaRuleList = res.data.metaDataList.filter(item => item.rule);
          this.metaDataList = res.data.metaDataList.filter(item => !item.rule);
        })
        .catch(() => {
          this.loading.itemStatus = false;
        });
    },
    handleClick(item, index) {
      this.dialogVisible = true;
      // this.$set(item, "id", index + 1);
      const obj = JSON.parse(JSON.stringify(item));
      this.$nextTick(() => {
        this.cardObj = [obj];
        this.activeStatus = obj.available;
      });
      // this.cardObj[0].id = ;
    },
    handleCurrentChange(val) {
      this.page.pageNo = val;
      this.getList();
    },
    handleSizeChange(val) {
      this.page.pageSize = val;
      this.getList();
    },
    handleDetailStatusChange(val) {
      this.loading.dialog = true;
      setSectionStatus({
        docId: this.obj.id,
        contentStatus: String(val),
        contentId: this.cardObj[0].contentId,
        all:false,
      })
        .then((res) => {
          this.loading.dialog = false;
          if (res.code === 0) {
            this.$message.success(this.$t('knowledgeManage.operateSuccess'));
          } else {
            this.cardObj[0].available = !this.cardObj[0].available;
          }
        })
        .catch(() => {
          this.loading.dialog = false;
          this.cardObj[0].contentStatus = !this.cardObj[0].contentStatus;
        });
    },
    handleStatusChange(item, index) {
      this.loading.itemStatus = true;
      setSectionStatus({
        docId: this.obj.id,
        contentStatus: String(item.available),
        contentId: item.contentId,
        all: false,
      })
        .then((res) => {
          this.loading.itemStatus = false;
          if (res.code === 0) {
            this.$message.success(this.$t('knowledgeManage.operateSuccess'));
            this.getList();
          } else {
            this.res.contentList[index].available =
              !this.res.contentList[index].available;
          }
        })
        .catch(() => {
          this.res[index].contentStatus = !this.res[index].contentStatus;
          this.loading.itemStatus = false;
        });
    },
    handleStatus(type) {
      this.loading.itemStatus = true;
      setSectionStatus({
        docId: this.obj.id,
        contentStatus: type==='start' ? "true" :"false",
        contentId: "",
        all:true,
      })
        .then((res) => {
          this.loading.itemStatus = false;
          if (res.code === 0) {
            this.$message.success(this.$t('knowledgeManage.operateSuccess'));
            this.getList();
          }
        })
        .catch(() => {
          this.loading.itemStatus = false;
        });
    },
    renderHeader(h, { column, $index }) {
      // column列数据 $index当前列索引
      const columnHtml =
        this.$t('knowledgeManage.section') +
        this.cardObj[0].contentNum +
        this.$t('knowledgeManage.length')+ " :" +
        this.cardObj[0].content.length +
        this.$t('knowledgeManage.character');
      return h("span", {
        domProps: {
          innerHTML: columnHtml,
        },
      });
    },
    handleClose() {
      this.dialogVisible = false;
      if (this.cardObj[0].available === this.activeStatus) return;
      this.getList();
    },
  },
};
</script>
<style lang="scss">
.showMore{
  margin-left:5px;
  background:#f4f5ff;
  padding:2px;
  border-radius:4px;
}
.metaItem{
  margin-left:5px;
  background:#f4f5ff;
  padding:2px;
  border-radius:4px;
}
.editIcon{
  cursor: pointer;
  color:#384BF7;
  font-size:16px;
  display: inline-block;
  margin-left:5px;
}
.section {
  width: 100%;
  height: 100%;
  padding: 20px 20px 30px 20px;
  margin: auto;
  overflow: auto;
  //background: #fafafa;

  .el-divider--horizontal {
    margin: 30px 0;
  }
  .title {
    font-size: 18px;
    font-weight: bold;
    color: #333;
    padding: 10px 0;
  }

  .container {
    min-width: 980px;
    padding: 15px;
    height: calc(100% - 45px);
    /*background: #fff;
    box-shadow: 0 1px 6px rgba(0, 0, 0, 0.3);*/
    border-radius: 5px;
    overflow: auto;

    .el-descriptions :not(.is-bordered) .el-descriptions-item__cell {
      &:nth-child(even) {
        width: 25%;
      }
      padding: 10px;
      // font-size: 12px;
    }
    .btn {
      padding: 10px 0;
      text-align: right;
    }

    .card {
      flex-wrap: wrap;
      // display: flex;
      // justify-content: space-between;

      .el-row {
        margin: 0 !important;
      }
      .text {
        font-size: 14px;
      }

      .item {
        height: 120px;
        margin-bottom: 18px;
        display: -webkit-box;
        -webkit-line-clamp: 6;
        -webkit-box-orient: vertical;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .clearfix:before,
      .clearfix:after {
        display: table;
        content: "";
      }
      .clearfix:after {
        clear: both;
      }

      .card-box {
        margin-bottom: 10px;

        .box-card {
          &:hover {
            cursor: pointer;
            transform: scale(1.03);
          }
        }
      }

      .el-card__header {
        padding: 8px 20px;
      }
    }
  }
}
</style>