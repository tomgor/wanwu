<template>
    <div class="page-wrapper full-content">
        <div class="page-title">
            <i class="el-icon-arrow-left" @click="goBack" style="margin-right: 10px; font-size: 20px; cursor: pointer">
            </i>
            {{$t('knowledgeManage.hitTest')}}
        </div>
        <div class="block wrap-fullheight">
            <div class="test-left test-box">
                <div class="hitTest_input">
                    <h3>命中分段测试</h3>
                    <el-input type="textarea" :rows="4" v-model="question" class="test_ipt"/>
                    <div class="test_btn">
                        <el-button type="primary" size="small" @click="startTest">开始测试<span class="el-icon-caret-right"></span></el-button>
                    </div>
                </div>
                <el-form :model="formInline" ref="formInline" :inline="false" class="test_form">
                    <el-form-item label="选择知识库" class="vertical-form-item">
                        <el-select v-model="formInline.knowledgeIdList" placeholder="请选择" multiple clearable style="width:100%;">
                            <el-option
                                v-for="item in knowledgeOptions"
                                :key="item.knowledgeId"
                                :label="item.name"
                                :value="item.knowledgeId">
                            </el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="Rerank模型" class="vertical-form-item">
                        <el-select 
                        clearable
                        style="width:100%;"
                        loading-text="模型加载中..."
                        v-model="formInline.rerankModelId" 
                        placeholder="请选择">
                            <el-option
                                v-for="item in rerankOptions"
                                :key="item.modelId"
                                :label="item.displayName"
                                :value="item.modelId">
                            </el-option>
                        </el-select>
                    </el-form-item>
                </el-form>
            </div>
            <div class="test-right test-box">
                <div class="result_title">
                    <h3>命中预测结果</h3>
                    <img src="@/assets/imgs/nodata_2x.png" v-if="searchList.length >0" />
                </div>
                <div class="result" v-loading="resultLoading">
                    <div v-if="searchList.length >0" class="result_box">
                        <div v-for="(item,index) in searchList" :key="'result'+index" class="resultItem">
                                <div class="resultTitle">
                                    <span class="tag">{{$t('knowledgeManage.section')}}{{index+1}}</span>
                                    <span class="score">{{$t('knowledgeManage.hitScore')}}: {{score[index]}}</span>
                                </div>
                                <div>
                                    <div>{{item.title + ', ' + item.snippet}}</div>
                                    <div class="file_name">文件名称：{{item.title}}</div>
                                </div>
                        </div>
                    </div>
                    <div v-else class="nodata">
                        <img src="@/assets/imgs/nodata_2x.png" />
                        <p class="nodata_tip">暂无数据</p>
                    </div>
                </div>
            </div>
        </div>
            <el-main class="padding border" style="margin-left:10px;">
                <el-container style="width:100%;height:100%">
                    <el-header>{{$t('knowledgeManage.hitPrediction')}}</el-header>
                    <el-main class="padding">
                       <div style="height:30%;width:100%;">
                            <el-input  type="textarea" :rows="4" :placeholder="$t('knowledgeManage.inputTestContent')" v-model="testInput"></el-input>
                            <div class="btn">
                                <el-button type="primary" @click="startTest">{{$t('knowledgeManage.startTest')}}</el-button>
                            </div>
                       </div>
                       <div style="height:70%;width:100%;">
                        <el-header>{{$t('knowledgeManage.hitResult')}}</el-header>
                        <div class="result" v-loading="resultLoading">
                            <template v-if="searchList.length >0">
                            <div v-for="(item,index) in searchList" :key="'result'+index" class="resultItem">
                                <div class="resultTitle">
                                    <span class="tag">{{$t('knowledgeManage.section')}}{{index+1}}</span>
                                    <span class="score">{{$t('knowledgeManage.hitScore')}}: {{score[index]}}</span>
                                </div>
                                <div>{{item}}</div>
                            </div>
                            </template>
                            <div v-else class="noResult">
                               {{ noResult }}
                            </div>
                        </div>
                       </div>
                    </el-main>
                </el-container>
            </el-main>
        </el-container> -->
    </div>
</template>
<script>
import { getKnowledgeList,hitTest } from "@/api/knowledge";
import { getRerankList} from "@/api/modelAccess";
export default{
    data(){
        return{
            knowledgeOptions:[],
            rerankOptions:[],
            formInline:{
                knowledgeIdList:[],
                rerankModelId:''
            },
            question:'',
            resultLoading:false,
            searchList:[
                { "title": "名胜古迹简介.docx", "snippet": "故宫简介： \n北京故宫（The Imperial Palace ），位于中国北京市，是明清两代的皇家宫殿，旧称紫禁城，位于北京中轴线的中心。以三大殿为中心，占地面积约72万平方米，建筑面积约15万平方米，有大小宫殿七十多座，相传故宫一共有9999.5间房屋，实际据1973年专家现场测量故宫有房间8707间 。 \n故宫于明永乐四年（1406年）开始建设，以南京故宫为蓝本营建，到永乐十八年（1420年）建成，成为明清两朝二十四位皇帝的皇宫。民国十四年（1925年）10月10日，故宫博物院正式成立开幕。北京故宫南北长961米，东西宽753米，四面围有高10米的城墙，城外有宽52米的护城河。故宫有四座城门，南面为午门，北面为神武门，东面为东华门，西面为西华门。城墙的四角各有一座风姿绰约的角楼，民间有九梁十八柱七十二条脊之说，形容其结构的复杂。 \n故宫内的建筑分为外朝和内廷两部分。外朝的中心为太和殿、中和殿、保和殿，统称三大殿，是国家举行大典礼的地方。三大殿左右两翼辅以文华殿、武英殿两组建筑。内廷的中心是乾清宫、交泰殿、坤宁宫，统称后三宫，是皇帝和皇后居住的正宫。其后为御花园。后三宫两侧排列着东、西六宫，是后妃们居住休息的地方。东六宫东侧是天穹宝殿等佛堂建筑，西六宫西侧是中正殿等佛堂建筑。外朝、内廷之外还有外东路、外西路两部分建筑。 \n故宫是世界上现存规模最大、保存最为完整的木质结构古建筑群之一。1961年3月4日，北京故宫被公布为第一批全国重点文物保护单位。1987年被列为世界文化遗产。 \n中国国家博物馆简介：", "knowledgeName": "名胜古迹介绍", "meta_data": {} },
                { "title": "名胜古迹简介123.docx", "snippet": "故宫简介： \n北京故宫（The Imperial Palace ），位于中国北京市，是明清两代的皇家宫殿，旧称紫禁城，位于北京中轴线的中心。以三大殿为中心，占地面积约72万平方米，建筑面积约15万平方米，有大小宫殿七十多座，相传故宫一共有9999.5间房屋，实际据1973年专家现场测量故宫有房间8707间 。 \n故宫于明永乐四年（1406年）开始建设，以南京故宫为蓝本营建，到永乐十八年（1420年）建成，成为明清两朝二十四位皇帝的皇宫。民国十四年（1925年）10月10日，故宫博物院正式成立开幕。北京故宫南北长961米，东西宽753米，四面围有高10米的城墙，城外有宽52米的护城河。故宫有四座城门，南面为午门，北面为神武门，东面为东华门，西面为西华门。城墙的四角各有一座风姿绰约的角楼，民间有九梁十八柱七十二条脊之说，形容其结构的复杂。 \n故宫内的建筑分为外朝和内廷两部分。外朝的中心为太和殿、中和殿、保和殿，统称三大殿，是国家举行大典礼的地方。三大殿左右两翼辅以文华殿、武英殿两组建筑。内廷的中心是乾清宫、交泰殿、坤宁宫，统称后三宫，是皇帝和皇后居住的正宫。其后为御花园。后三宫两侧排列着东、西六宫，是后妃们居住休息的地方。东六宫东侧是天穹宝殿等佛堂建筑，西六宫西侧是中正殿等佛堂建筑。外朝、内廷之外还有外东路、外西路两部分建筑。 \n故宫是世界上现存规模最大、保存最为完整的木质结构古建筑群之一。1961年3月4日，北京故宫被公布为第一批全国重点文物保护单位。1987年被列为世界文化遗产。 \n中国国家博物馆简介：", "knowledgeName": "名胜古迹介绍", "meta_data": {} }
            ],
            score:[0.5,0.6],
            noResult: "",
        }
    },
    created(){
        this.getKnowledgeList();
        this.getRerankData();
    },
    methods:{
        async getKnowledgeList() {
            //获取文档知识分类
            const res = await getKnowledgeList({});
            if (res.code === 0) {
                this.knowledgeOptions = res.data.knowledgeList || [];
            } else {
                this.$message.error(res.message);
            }
        },
        getRerankData(){
            getRerankList().then(res =>{
                if(res.code === 0){
                    this.rerankOptions = res.data.list || []
                }
            })
        },
        goBack(){
          this.$router.go(-1);
        },
        startTest(){
            if(this.formInline.knowledgeIdList.length === 0 ){
                this.$message.warning(this.$t('knowledgeManage.pselectKnowledgeTips'))
                return
            }
            if(this.formInline.rerankModelId.length === 0 ){
                this.$message.warning(this.$t('knowledgeManage.pselectKnowledgeTips'))
                return
            }
            if(this.question === ''){
                 this.$message.warning('请选择Rerank模型')
                return
            }
            const data = {
                ...this.formInline,
                question:this.question
            }
            this.test(data)
        },
        async test(data){
            this.resultLoading = true
            const res = await hitTest(data);
            if(res.code === 0){
                this.$message.success(this.$t('knowledgeManage.operateSuccess'))
                this.searchList = res.data !== null ? res.data.searchList : [];
                if(res.data){
                    this.score = res.data.score.map(item =>item.toFixed(5))
                }else{
                    this.noResult = this.$t('knowledgeManage.testResultTips')+`“${data.question}”`+this.$t('knowledgeManage.testResultTips1');
                    this.score = []
                }
                this.resultLoading = false
            }else{
                this.searchList = [];
                this.resultLoading = false
            }
        }
    }
}
</script>
<style lang="scss" scoped>
/deep/{
    .vertical-form-item{
        display: flex;
        flex-direction: column;
        align-items: flex-start;
    }
    .vertical-form-item .el-form-item__label {
        line-height:unset;
        font-size:14px;
        font-weight:bold;
    }
    .el-form-item__content{
        width:100%;
    }
}
.full-content{
    display:flex;
    flex-direction: column;
    .page-title{
        border-bottom:1px solid #d9d9d9;
    }
    .block{
        margin:30px 10px;
        display:flex;
        gap:20px;
      .test-box{
        flex:1;
        height:100%;
        .hitTest_input{
            background:#fff;
            border-radius:6px;
            border:1px solid #e9ecef;
            padding:0 20px;
            h3{
                padding:30px 0 10px 0;
                font-size:14px;
                font-weight:bold;
            }
            .test_ipt{
                padding-bottom:10px;
            }
            .test_btn{
               padding:10px 0;
               display:flex;
               justify-content:flex-end;
            }
        }
        .test_form{
            margin-top:20px;
            padding:20px;
            background:#fff;
            border-radius:6px;
            border:1px solid #e9ecef;
        }
      }
      .test-right{
            background:#fff;
            border-radius:6px;
            border:1px solid #e9ecef;
            height:100%;
            padding:20px;
            display:flex;
            flex-direction: column;
            .result_title{
                display:flex;
                justify-content:space-between;
                h3{
                    padding: 10px 0 10px 0;
                    font-size: 14px;
                }
                img{
                    width:150px;
                }
            }
           
            .result{
                flex: 1;
                width:100%;
                .result_box{
                    width:100%;
                    height:100%;
                    overflow-y:auto;
                    .resultItem{
                        background:#F7F8FA; 
                        border-radius:6px;
                        margin-bottom:20px;
                        padding:20px;
                        color:#666666;
                        line-height:1.8;
                        .resultTitle{
                            display:flex;
                            align-items:center;
                            justify-content:space-between;
                            padding: 10px 0;
                        .tag{
                            color: #384BF7;
                            display:inline-block;
                            background: #D2D7FF ;
                            padding:0 10px;
                            border-radius:6px;
                        }
                        .score{
                            color: #384BF7;
                            font-weight:bold;
                        }
                        }
                        .file_name{
                            border-top:1px dashed #d9d9d9;
                            margin:10px 0;
                            padding-top:10px;
                            font-weight:bold;
                        }
                  }
                }
                .nodata{
                    width:100%;
                    height:100%;
                    display:flex;
                    align-items:center;
                    justify-content:center;
                    flex-direction: column;
                     align-self: center; /* 仅该元素纵向居中 */
                    .nodata_tip{
                        padding:10px 0;
                        color: #595959 ;
                    }
                }
            }
            
        }
    }
    
}

</style>