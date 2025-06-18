<template>
    <div class="mcp-detail">
        <span class="back" @click="back">返回MCP广场</span>
        <div class="mcp-title">
            <svg-icon icon-class="mcp_server" class="editable--send logo" />
            <div :class="['info',{fold:foldStatus}]">
                <p class="name">{{detail.name}}</p>
                <p v-if="detail.description.length>260" class="desc">{{foldStatus?detail.description:detail.description.slice(0,268)+'...'}}
                    <span class="arrow" v-show="detail.description.length>260" @click="fold"> {{foldStatus?'收起':'详情 >>'}}</span>
                </p>
                <p v-else class="desc">{{detail.description}}</p>
            </div>
        </div>
        <div class="main">
            <div class="left-info">
                <!-- tabs -->
                <div class="mcp-tabs">
                    <div v-if="source === '2'" :class="['mcp-tab',{ 'active': tabActive === 0 }]" @click="tabClick(0)">介绍概览</div>
                    <div style="display: inline-block">
                        <div :class="['mcp-tab',{ 'active': tabActive === 1 }]" @click="tabClick(1)">SSE URL及工具</div>
                    </div>
                </div>

                <div v-if="tabActive === 0">
                    <div class="overview bg-border" >
                        <div class="overview-item">
                            <div class="item-title">• &nbsp;使用概述</div>
                            <div class="item-desc" v-html="parseTxt(squareDetail.whatIs)"></div>
                        </div>
                        <div class="overview-item">
                            <div class="item-title">• &nbsp;特性说明</div>
                            <div class="item-desc" v-html="parseTxt(squareDetail.keyFeatures)"></div>
                        </div>
                        <div class="overview-item">
                            <div class="item-title">• &nbsp;应用场景</div>
                            <div class="item-desc" >
                                <div v-html="parseTxt(squareDetail.whereToUse)"></div>
                                <div v-html="parseTxt(squareDetail.useCases)"></div>
                            </div>
                        </div>
                    </div>
                    <div class="overview bg-border" >
                        <div class="overview-item">
                            <div class="item-title">• &nbsp;使用说明</div>
                            <div class="item-desc" v-html="parseTxt(squareDetail.howToUse)"></div>
                        </div>
                    </div>
                    <div class="overview bg-border" >
                        <div class="overview-item">
                            <div class="item-title">• &nbsp;详情</div>
                            <div class="item-desc">
                                <div class="readme-content markdown-body mcp-markdown" v-html="md.render(markdownHtml)"></div>
                            </div>

                        </div>
                    </div>
                </div>
                <div v-if="tabActive === 1">
                    <div class="tool bg-border">
                        <div class="tool-item ">
                            <p class="title">SSE URL:</p>
                            <div class="sse-url" style="display: flex">
                                <div class="tool-item-bg sse-url__input">{{detail.serverUrl}}</div>
                            </div>
                            <p style="line-height: 40px;color: #666;">* 您已添加到自定义，可直接在工作流或智能体中直接调用。</p>
                        </div>
                        <div class="tool-item">
                            <p class="title">工具介绍:</p>
                            <div class="tool-item-bg tool-intro">
                                <el-collapse>
                                    <el-collapse-item v-for="(n,i) in tools" :title="n.name" :name="i" :key="n.name + i">
                                        <div class="desc">描述：{{n.description}}</div>
                                        <div class="params">
                                            <p>参数说明:</p>
                                            <div class="params-table" v-for="(m,j) in n.params">
                                                <div class="tr">
                                                    <div class="td">{{m.name}}</div>
                                                    <div class="td color">{{m.type}}</div>
                                                    <div class="td color">{{m.requiredBadge}}</div>
                                                </div>
                                                <p class="params-desc">{{m.description}}</p>
                                            </div>
                                        </div>
                                    </el-collapse-item>
                                </el-collapse>
                            </div>
                        </div>
                        <!--<div class="tool-item">
                            <p class="title">MCP服务器配置:</p>
                            <div class="tool-item-bg service-config"></div>
                        </div>-->
                        <div class="tool-item">
                            <p class="title">安装说明:</p>
                            <div class="tool-item-bg">
                                <div class="install-intro-item">
                                    <p class="install-intro-title">在 Cursor 中安装</p>
                                    <p>1. 点击Cursor右上角'设置'，进入左侧菜单中的'MCP'选项</p>
                                    <p>2. 点击页面右上方的'+添加'按钮, 自动打开mcp.json配置文件</p>
                                    <p>3. 在文件中粘贴MCP配置并保存（在合适位置插入，无需删除已有内容）</p>
                                    <p>4. MCP设置界面显示绿色原点即可使用</p>
                                </div>
                                <div class="install-intro-item">
                                    <p class="install-intro-title">在 Claude 中安装</p>
                                    <p>1. 在Claude页面左上角打开'设置'，进入'Developer'</p>
                                    <p>2. 点击'Edit Config' 定位claude_desktop_config.json配置文件</p>
                                    <p>3. 在json文件中粘贴MCP配置并保存（在合适位置插入，无需删除已有内容）</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>


            </div>

            <div class="right-recommend">
                <p style="margin: 20px 0;color: #333;">其他MCP服务查看</p>
                <div class="recommend-item" v-for="(n,i) in recommendList" @click="handleClick(n)">
                    <img class="logo" :src="n.logo"/>
                    <p class="name">{{n.name}}</p>
                    <p class="intro">{{n.description}}</p>
                </div>
            </div>
        </div>

    </div>
</template>
<script>
    import { getDetail, getTools, getPublicMcpList,getMarkDownContent } from "@/api/mcp";
    import {md} from '../../mixins/marksown-it'

    export default {
        data() {
            return {
                md:md,
                mcpId:'',
                mcpSquareId:'',
                source:'',
                squareDetail: {
                    description:'',
                    whatIs:'',
                    keyFeatures:'',
                    whereToUse:'',
                    useCases:'',
                    howToUse:'',
                },
                markdownHtml:'',
                detail:{
                    description:'',
                },
                tools:[],
                foldStatus:false,
                tabActive:1,
                recommendList:[],
            };
        },
        created() {
            this.initData()
        },
        methods: {
            back(){
              this.$router.go(-1)
            },
            tabClick(status){
                this.tabActive = status
            },
            initData(){
                this.mcpId = this.$route.params.mcpId
                this.source = this.$route.params.source

                console.log(this.mcpId, this.source)
                this.getDataFromCustomize()
                this.getRecommendList()

                if(this.source === '2'){
                    this.tabActive = 0;
                    this.mcpSquareId = this.$route.params.mcpSquareId
                    this.getDataFromSquare()
                }
            },
            //===来自广场的详情===
            getDataFromSquare(){
                getPublicMcpList({mcpSquareId:this.mcpSquareId,pageNo: 1, pageSize: 1})
                    .then((res) => {
                        this.squareDetail = res.data.list[0];
                        this.loadMarkDown()
                    })
                    .catch((err) => {});
            },
            //读取markdown文档
            loadMarkDown(){
                getMarkDownContent(this.squareDetail.contentPath.split('https://obs-nmhhht6.cucloud.cn/')[1])
                    .then(async(res) => {
                        const text = await res.text();
                        this.markdownHtml = text
                    })
                    .catch((err) => {});
            },
            //解析文本，遇到.换行等
            parseTxt(txt){
                return txt.replaceAll('\n\t','<br/>&nbsp;').replaceAll('\t', '   &nbsp;')
            },
            //===自定义详情===
            getDataFromCustomize(){
                getDetail({mcpId:this.mcpId})
                    .then((res) => {
                        this.detail = res.data;
                        this.doGetTools()
                    })
                    .catch((err) => {});
            },
            doGetTools(){
                getTools({
                    serverUrl: this.detail.serverUrl,
                })
                    .then((res) => {
                        //解析tools 把object转成array
                        this.tools = res.data.tools.map((n,i)=>{
                            let params = []
                            let properties = n.inputSchema.properties
                            for(var key in properties){
                                params.push({
                                    "name": key,
                                    "requiredBadge": n.inputSchema.required.includes(key)?'必填':'',
                                    "type": properties[key].type,
                                    "description": properties[key].description,
                                })
                            }
                            return {
                                ...n,
                                params
                            }
                        })
                    })
                    .catch((err) => {
                    });
            },
            getRecommendList(){
                let params = {
                    pageNo: 1,
                    pageSize: 10,
                    hosted:false
                }
                getPublicMcpList(params)
                    .then((res) => {
                        this.recommendList = res.data.list;
                    })
                    .catch((err) => {});
            },
            handleClick(val){
                this.$router.push(`/mcp/public/detail/${val.mcpSquareId}/${val.hosted?1:0}`)
            },
            fold(){
                this.foldStatus = !this.foldStatus
            },

        },
    };
</script>
<style lang="scss">
@import "../mcpManagementPublic/markdown.min.css";
.markdown-body{
  font-family: 'Microsoft YaHei', Arial, sans-serif;
  color: #333;
}
.mcp-detail{
  padding: 20px;
  /*height: calc(100vh - 88px);*/
  overflow: auto;
  margin: 20px;
  /*background: linear-gradient(180deg, hsla(0, 0%, 94.5%, .83), hsla(0, 0%, 94.5%, .1) 91%);
  box-shadow: 0 3px 8px 0 rgba(0, 0, 0, .1);*/
  .back {
    color: $color;
    cursor: pointer;
  }
  .mcp-title{
    padding: 20px 0;
    display: flex;
    border-bottom: 1px solid #bfbfbf;
    .logo{
      width: 54px;
      height: 54px;
      object-fit: cover;
    }
    .info{
      position: relative;
      width: 1240px;
      margin-left: 15px;
      .name{
        font-size: 16px;
        color: #5D5D5D;
        font-weight: bold;
      }
      .desc{
        margin-top: 10px;
        line-height: 22px;
        color: #9F9F9F;
        word-break: break-all;
      }
      .arrow{
        position: absolute;
        display: block;
        right: 0;
        bottom: -5px;
        cursor: pointer;
        color: $color;
        margin-left: 10px;
        font-size: 13px;
      }
    }
    .fold{
      height: auto;
    }
  }
}
.main{
  display: flex;
  margin: 10px 0 50px 0;
  .left-info{
    width: calc(100% - 420px);
    margin-right: 20px;
    .mcp-tabs{
      margin: 20px 0 0 0;
      .mcp-tab{
        display: inline-block;
        vertical-align: middle;
        width: 160px;
        height: 40px;
        border-bottom: 1px solid #333;
        line-height: 40px;
        text-align: center;
        cursor: pointer;
      }
      .active{
        background: #333;
        color: #fff;
        font-weight: bold;
      }
    }
    .overview{
      .overview-item{
        display: flex;
        padding: 15px 0;
        border-bottom: 1px solid #eee;
        line-height: 24px;
        .item-title{
          width: 80px;
          color: $color;
          font-weight: bold;
        }
        .item-desc{
          width: calc(100% - 100px);
          margin-left: 10px;
          flex:1;
          color: #333;
        }

      }
      .overview-item:last-child{
        border-bottom: none;
      }
    }
    .markdown{

    }
    .install{

    }
    .tool{
      .tool-item{
        padding: 20px 0;
        border-bottom: 1px solid #eee;
        .title{
          font-weight: bold;
          line-height: 46px;
        }
        .tool-item-bg{
          background: inherit;
          background-color: rgba(249, 249, 249, 1);
          border: none;
          border-radius: 10px;
          padding: 20px;
        }
      }
      .tool-item:last-child{
        border-bottom: none;
      }
      .sse-url{
        .sse-url__input{
          width: 100%;
          padding: 12px;
          color: $color;
        }
        .sse-url__bt{
          width: 120px;
        }
      }
      .install-intro-item{
        p{
          line-height: 26px;
          color: #333;
        }
        .install-intro-title{
          color: $color;
          margin-top: 10px;
          font-weight: bold;
        }
      }
    }
  }
  .right-recommend{
    width: 400px;
    overflow-y: auto;
    border-left:1px solid #eee;
    padding: 20px;
    .recommend-item{
      position: relative;
      border: 1px solid $border_color; // rgba(208, 167, 167, 1);
      background: #F4F5FF; // rgba(255, 247, 247, 1);
      margin-bottom: 15px;
      border-radius: 10px;
      padding: 20px 20px 20px 80px;
      text-align: left;
      .logo{
        width: 46px;
        height: 46px;
        object-fit: cover;
        position: absolute;
        left:20px;
        border: 1px solid #fff;
        border-radius: 4px;
      }
      .name{
        color: #5D5D5D;
        font-weight: bold;
      }
      .intro{
        height: 34px;
        color: #5D5D5D;
        margin-top: 8px;
        font-size: 13px;
        overflow: hidden;
      }
    }
  }
}
.bg-border{
  margin-top: 20px;
  /*min-height: calc(100vh - 360px);*/
  background-color: rgba(255, 255, 255, 1);
  box-sizing: border-box;
  /*border:1px solid rgba(208, 167, 167, 1);*/
  border-radius: 10px;
  padding: 10px 20px;
  box-shadow: 2px 2px 15px #d0a7a757;
}
.overview-item .item-desc{
  line-height: 28px;
}
.mcp-el-collapse.el-collapse{
  border: none;
}
.mcp-el-collapse .el-collapse-item{
  margin: 10px 0;
  border: none;
  .el-collapse-item__header{
    border: none;
    color: $color;
    font-weight: bold;
    padding: 0 20px;
  }
  .el-collapse-item__wrap{
    padding: 0 20px;
    border: none;
  }
  .desc{
    background: rgba(255, 246, 246, 1);
    padding: 10px 15px;
    border-radius: 6px;
    border: 1px solid #f5cbcb;
  }
  .params{
    margin-top: 12px;
    .params-table{
      border-radius: 6px;
      border: 1px solid #ddd;
      padding: 10px 12px;
      background-color: #fff;
      margin-top: 6px;
      .tr{
        display: flex;
        .td{
          padding: 0 30px 0 0;
        }
        .color{
          color: $color;
        }
      }
      .params-desc{
        margin-top: 4px;
        color: #999;
      }
    }
  }
}
</style>