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
          <el-form-item
            label="检索方式配置"
            class="vertical-form-item"
          >
          <div v-for="item in searchTypeData" :class="['searchType-list',{ 'active': item.showContent }]">
            <div class="searchType-title" @click="clickSearch(item)">
              <span :class="[item.icon,'img']"></span>
              <div class="title-content">
                <div class="title-box">
                  <h3 class="title-name">{{item.name}}</h3>
                  <p class="title-desc">{{item.desc}}</p>
                </div>
                <span :class="item.showContent?'el-icon-arrow-up':'el-icon-arrow-down'"></span>
              </div>
            </div>
            <div class="searchType-content" v-if="item.showContent">
              <div v-if="item.isWeight" class="weightType-box" >
                <div v-for="mixItem in item.mixType" :class="['weightType',{ 'active': mixItem.value === item.mixTypeValue }]" @click.stop="mixTypeClick(item,mixItem)">
                  <p class="weightType-name">{{mixItem.name}}</p>
                  <p class="weightType-desc">{{mixItem.desc}}</p>
                </div>
              </div>
              <el-row v-if="item.isWeight && item.mixTypeValue === 'weight'" @click.stop>
                <el-col class="mixTypeRange-title">
                  <span>语义[{{item.mixTypeRange[0]}}]</span>
                  <span>关键词[{{item.mixTypeRange[1]}}]</span>
                </el-col>
                <el-col>
                  <el-slider
                  v-model="item.mixTypeRange"
                  range
                  show-stops
                  :step="0.1"
                  :max="1"
                  >
                </el-slider>
                </el-col>
              </el-row>
              <el-row v-if="showRerank(item)">
                <el-col class="content-name">Rerank模型</el-col>
                <el-col>
                  <el-select
                  clearable
                  filterable 
                  style="width:100%;"
                  loading-text="模型加载中..."
                  v-model="formInline.knowledgeMatchParams.rerankModelId"
                  @visible-change="visibleChange($event)"
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
                </el-col>
              </el-row>
              <el-row>
                <el-col>
                  <span class="content-name">TopK</span>
                  <el-tooltip class="item" effect="dark" content="用于控制检索阶段返回的最相关的文档片段的数量。这些文档片段将被送入生成模型中，用于 生成最终的回答。" placement="right">
                    <span class="el-icon-question tips"></span>
                  </el-tooltip>
                </el-col>
                <el-col>
                  <el-slider
                    :min="1"
                    :max="10"
                    :step="1"
                    v-model="formInline.knowledgeMatchParams.topK"
                    show-input>
                  </el-slider>
                </el-col>
              </el-row>
              <el-row>
                <el-col>
                  <span class="content-name">Score阈值</span>
                   <el-tooltip class="item" effect="dark" content="检索结果的相似度阈值，低于该值的结果将被过滤。" placement="right">
                      <span class="el-icon-question tips"></span>
                   </el-tooltip>
                </el-col>
                <el-col>
                  <el-slider
                    :min="0"
                    :max="1"
                    :step="0.1"
                    v-model="formInline.knowledgeMatchParams.score"
                    show-input>
                  </el-slider>
                </el-col>
              </el-row>
            </div>
          </div>
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
import { hitTest } from "@/api/knowledge";
import { getRerankList } from "@/api/modelAccess";
import { md } from "@/mixins/marksown-it";
export default {
  data() {
    return {
      md: md,
      rerankOptions: [],
      formInline: {
        knowledgeIdList: [this.$route.query.knowledgeId],
        knowledgeMatchParams:{
          keywordPriority:0.8,//关键词权重
          matchType:'',//vector（向量检索）、text（文本检索）、mix（混合检索：向量+文本）
          priorityMatch:1,//权重匹配，只有在混合检索模式下，选择权重设置后，这个才设置为1
          rerankModelId:'',//rerank模型id
          score:0.4,//过滤分数阈值
          semanticsPriority:0.2,//语义权重
          topK:1//topK 获取最高的几行
        }
      },
      question: "",
      resultLoading: false,
      searchList: [],
      score: [],
      searchTypeData:[
        {
          name:'向量检索',
          value:'vector',
          desc:'通过向量相似度找到语义相近、表达多样的文本片段，适用于理解和召回语义相关信息。',
          rerank:'',
          icon:'el-icon-menu',
          topK:0,
          Score:0.4,
          isWeight:false,
          showContent:false
        },
        {
          name:'全文检索',
          value:'text',
          desc:'基于关键词匹配，能够高效查询包含指定词汇的文本片段，适用于精确查找',
          rerank:'',
          topK:0,
          Score:0.4,
          icon:'el-icon-document',
          isWeight:false,
          showContent:false
        },
        {
          name:'混合检索',
          value:'mix',
          desc:'结合向量和关键词检索，融合语义理解与关键词匹配，兼顾相关性和准确性，提升检索效果。',
          rerank:'',
          icon:'el-icon-s-grid',
          topK:0,
          Score:0.4,
          isWeight:true,
          Weight:'',
          mixTypeValue:'weight',
          showContent:false,
          mixTypeRange:[0.2,0.8],
          mixType:[
            {
              name:'权重设置',
              value:'weight',
              desc:'权重设置功能用于调整不同检索方式的影响力。通过设置权重，可以控制语义相似度和关键词匹配在最终排序中的占比。'
            },
            {
              name:'Rerank模型',
              value:'rerank',
              desc:'重排序模型会根据候选文档与用户问题的语义匹配度，对初步检索结果进行重新排序从而进一步提升最终返回结果的相关性和准确性。'
            }
            ]
        }
      ]
    };
  },
  created() {
    this.getRerankData();
  },
  methods: {
    mixTypeClick(item,n){
      item.mixTypeValue = n.value;
      this.formInline.knowledgeMatchParams.priorityMatch = n.value === 'weight' ? 1 : 0 ; 
    },
    showRerank(n){
      return (n.value === 'vector' || n.value === 'text') || (n.value === 'mix' && n.mixTypeValue === 'rerank');
    },
    clickSearch(n){
      this.formInline.knowledgeMatchParams.matchType = n.value;
      this.searchTypeData = this.searchTypeData.map(item => ({
        ...item,
        showContent: item.value === n.value ? !item.showContent : false
      }));
      this.clear();
    },
    clear(){
      this.formInline.knowledgeMatchParams.rerankModelId = '';
      this.formInline.knowledgeMatchParams.keywordPriority = 0.8;
      this.formInline.knowledgeMatchParams.semanticsPriority = 0.2;
      this.formInline.knowledgeMatchParams.priorityMatch = 1;
      this.formInline.knowledgeMatchParams.score = 0.4;
      this.formInline.knowledgeMatchParams.topK = 1;
    },
    getRerankData() {
      getRerankList().then((res) => {
        if (res.code === 0) {
          this.rerankOptions = res.data.list || [];
        }
      });
    },
    visibleChange(val) {
      if (val) {
        this.getRerankData();
      }
    },
    goBack() {
      this.$router.go(-1);
    },
    startTest() {
      const { matchType, priorityMatch, rerankModelId } = this.formInline.knowledgeMatchParams;
      if (this.question === "") {
        this.$message.warning("请输入问题");
        return;
      }

      if(matchType === ''){
        this.$message.warning("请选择检索方式");
        return;
      }
      if(matchType === 'mix' && priorityMatch === 1){
        this.formInline.knowledgeMatchParams.keywordPriority = this.searchTypeData[2]['mixTypeRange'][1];
        this.formInline.knowledgeMatchParams.semanticsPriority = this.searchTypeData[2]['mixTypeRange'][0];
      }else{
        if(rerankModelId === '') {
          this.$message.warning("请选择Rerank模型");
          return;
        }
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
.active{
  border:1px solid #384bf7!important;
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
       overflow-y:auto;
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
        .searchType-list:hover{
          border:1px solid #384bf7;
        }
        .searchType-list{
          border:1px solid #C0C4CC;
          border-radius:4px;
          margin:20px 0;
          padding:0 10px;
          cursor: pointer;
          .searchType-title{
            display:flex;
            align-items: center;
            .img{
              font-size:30px;
              text-align:center;
              line-height:50px;
              color: #384bf7;
              background-color:#fff;
              width:50px;
              height:50px;
              border-radius:8px;
              border:1px solid #e9e9eb;
              box-shadow:4px 2px 4px #f1f1f1 ;
            }
            .title-content{
              flex:1;
              display:flex;
              margin-left:10px;
              justify-content:space-between;
              align-items:center;
              .title-name{
                font-size: 16px;
                font-weight: bold;
                line-height: 1;
                padding-top:10px;
              }
              .title-desc{
                color:#888;
              }
            }
          }
          .searchType-content{
            padding:20px;
            .tips{
              color:#888;
              margin-left:5px;
            }
            .content-name{
              font-weight:bold;
            }
            .weightType-box{
              display:flex;
              gap: 20px;
              .weightType{
                border:1px solid #C0C4CC;
                border-radius:4px;
                .weightType-name{
                  text-align:center;
                  font-weight:bold;
                  line-height: 2;
                  font-size: 16px;
                  padding-top:5px;
                }
                .weightType-desc{
                  text-align:center;
                  line-height:1.5;
                  padding: 10px;
                  color:#888;
                }
              }
            }
            .mixTypeRange-title{
              display:flex;
              align-items:center;
              justify-content:space-between;
              font-weight:bold;
              margin-top:20px;
              line-height:1;
            }
          }
        }
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