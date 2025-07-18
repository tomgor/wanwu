<template>
  <div class="page-wrapper full-content">
    <div class="page-title">
      <i
        class="el-icon-arrow-left"
        @click="goBack"
        style="margin-right: 10px; font-size: 20px; cursor: pointer"
      >
      </i>
      {{$t('knowledgeManage.hitTest')}}
    </div>
    <div class="block wrap-fullheight">
      <div class="test-left test-box">
        <div class="hitTest_input">
          <h3>命中分段测试</h3>
          <el-input
            type="textarea"
            :rows="4"
            v-model="question"
            class="test_ipt"
          />
          <div class="test_btn">
            <el-button
              type="primary"
              size="small"
              @click="startTest"
            >开始测试<span class="el-icon-caret-right"></span></el-button>
          </div>
        </div>
        <el-form
          :model="formInline"
          ref="formInline"
          :inline="false"
          class="test_form"
        >
          <!-- <el-form-item
            label="选择知识库"
            class="vertical-form-item"
          >
            <el-select
              v-model="formInline.knowledgeIdList"
              placeholder="请选择"
              multiple
              clearable
              filterable 
              style="width:100%;"
              @visible-change="visibleChange($event,'knowledge')"
            >
              <el-option
                v-for="item in knowledgeOptions"
                :key="item.knowledgeId"
                :label="item.name"
                :value="item.knowledgeId"
              >
              </el-option>
            </el-select>
          </el-form-item> -->
          <el-form-item
            label="Rerank模型"
            class="vertical-form-item"
          >
            <el-select
              clearable
              filterable 
              style="width:100%;"
              loading-text="模型加载中..."
              v-model="formInline.rerankModelId"
              @visible-change="visibleChange($event,'rerank')"
              placeholder="请选择"
            >
              <el-option
                v-for="item in rerankOptions"
                :key="item.modelId"
                :label="item.displayName"
                :value="item.modelId"
              >
              </el-option>
            </el-select>
          </el-form-item>
        </el-form>
      </div>
      <div class="test-right test-box">
        <div class="result_title">
          <h3>命中预测结果</h3>
          <img
            src="@/assets/imgs/nodata_2x.png"
            v-if="searchList.length >0"
          />
        </div>
        <div
          class="result"
          v-loading="resultLoading"
        >
          <div
            v-if="searchList.length >0"
            class="result_box"
          >
            <div
              v-for="(item,index) in searchList"
              :key="'result'+index"
              class="resultItem"
            >
              <div class="resultTitle">
                <span class="tag">{{$t('knowledgeManage.section')}}{{index+1}}</span>
                <span class="score">{{$t('knowledgeManage.hitScore')}}: {{score[index]}}</span>
              </div>
              <div>
                <div v-html="md.render(item.snippet)" class="resultContent"></div>
                <div class="file_name">文件名称：{{item.title}}</div>
              </div>
            </div>
          </div>
          <div
            v-else
            class="nodata"
          >
            <img src="@/assets/imgs/nodata_2x.png" />
            <p class="nodata_tip">暂无数据</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import { getKnowledgeList, hitTest } from "@/api/knowledge";
import { getRerankList } from "@/api/modelAccess";
import { md } from "@/mixins/marksown-it";
export default {
  data() {
    return {
      md: md,
      knowledgeOptions: [],
      rerankOptions: [],
      formInline: {
        knowledgeIdList: [this.$route.query.knowledgeId],
        rerankModelId: "",
      },
      question: "",
      resultLoading: false,
      searchList: [],
      score: [],
    };
  },
  created() {
    // this.getKnowledgeList();
    this.getRerankData();
  },
  methods: {
    async getKnowledgeList() {
      //获取文档知识分类
      const res = await getKnowledgeList({});
      if (res.code === 0) {
        this.knowledgeOptions = res.data.knowledgeList || [];
      } else {
        this.$message.error(res.message);
      }
    },
    getRerankData() {
      getRerankList().then((res) => {
        if (res.code === 0) {
          this.rerankOptions = res.data.list || [];
        }
      });
    },
    visibleChange(val, type) {
      if (val) {
        if (type === "knowledge") {
          this.getKnowledgeList();
        } else {
          this.getRerankData();
        }
      }
    },
    goBack() {
      this.$router.go(-1);
    },
    startTest() {
      if (this.question === "") {
        this.$message.warning("请输入问题");
        return;
      }
      // if (this.formInline.knowledgeIdList.length === 0) {
      //   this.$message.warning(this.$t("knowledgeManage.pselectKnowledgeTips"));
      //   return;
      // }
      if (this.formInline.rerankModelId.length === 0) {
        this.$message.warning("请选择Rerank模型");
        return;
      }
      const data = {
        ...this.formInline,
        question: this.question,
      };
      this.test(data);
    },
    test(data) {
      this.resultLoading = true;
      this.searchList = [];
      this.score = [];
      hitTest(data).then((res) => {
        if (res.code === 0) {
          this.searchList = res.data !== null ? res.data.searchList : [];
          if (res.data) {
            this.score = res.data.score.map((item) => item.toFixed(5));
          } else {
            this.score = [];
          }
          this.resultLoading = false;
        } else {
          this.searchList = [];
          this.resultLoading = false;
        }
      }).catch(() =>{
        this.resultLoading = false;
      })
    },
  },
};
</script>
<style lang="scss" scoped>
/deep/ {
  .vertical-form-item {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
  }
  .vertical-form-item .el-form-item__label {
    line-height: unset;
    font-size: 14px;
    font-weight: bold;
  }
  .el-form-item__content {
    width: 100%;
  }
}
.full-content {
  display: flex;
  flex-direction: column;
  .page-title {
    border-bottom: 1px solid #d9d9d9;
  }
  .block {
    margin: 30px 10px;
    display: flex;
    height: calc(100% - 123px);
    gap: 20px;
    .test-box {
       flex:1;
       height:100%;
      .hitTest_input {
        background: #fff;
        border-radius: 6px;
        border: 1px solid #e9ecef;
        padding: 0 20px;
        h3 {
          padding: 30px 0 10px 0;
          font-size: 14px;
          font-weight: bold;
        }
        .test_ipt {
          padding-bottom: 10px;
        }
        .test_btn {
          padding: 10px 0;
          display: flex;
          justify-content: flex-end;
        }
      }
      .test_form {
        margin-top: 20px;
        padding: 20px;
        background: #fff;
        border-radius: 6px;
        border: 1px solid #e9ecef;
      }
    }
    .test-right {
      background: #fff;
      border-radius: 6px;
      border: 1px solid #e9ecef;
      height: 100%;
      padding: 20px;
      box-sizing:border-box;
      display: flex;
      flex-direction: column;
      .result_title {
        display: flex;
        justify-content: space-between;
        h3 {
          padding: 10px 0 10px 0;
          font-size: 14px;
        }
        img {
          width: 150px;
        }
      }
      .resultContent{
        img{
          width:100%;
        }
      }
      .result {
        flex: 1;
        width: 100%;
        display: flex;
        flex-direction: column;
        min-height: 0;
        .result_box {
          width: 100%;
          flex: 1;
          overflow-y: scroll;
          .resultItem {
            background: #f7f8fa;
            border-radius: 6px;
            margin-bottom: 20px;
            padding: 20px;
            color: #666666;
            line-height: 1.8;
            .resultTitle {
              display: flex;
              align-items: center;
              justify-content: space-between;
              padding: 10px 0;
              .tag {
                color: #384bf7;
                display: inline-block;
                background: #d2d7ff;
                padding: 0 10px;
                border-radius: 6px;
              }
              .score {
                color: #384bf7;
                font-weight: bold;
              }
            }
            .file_name {
              border-top: 1px dashed #d9d9d9;
              margin: 10px 0;
              padding-top: 10px;
              font-weight: bold;
            }
          }
        }
        .nodata {
          width: 100%;
          height: 100%;
          display: flex;
          align-items: center;
          justify-content: center;
          flex-direction: column;
          align-self: center; /* 仅该元素纵向居中 */
          .nodata_tip {
            padding: 10px 0;
            color: #595959;
          }
        }
      }
    }
  }
}
</style>