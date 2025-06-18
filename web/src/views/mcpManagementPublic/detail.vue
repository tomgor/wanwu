<template>
    <div class="mcp-detail"  id="timeScroll">
        <span class="back" @click="back">返回MCP广场</span>
        <div class="mcp-title">
            <img class="logo" :src="`${detail.logo}?x-oss-process=image/resize,w_100`"/>
            <div :class="['info',{fold:foldStatus}]">
                <p class="name">{{detail.alias || detail.name}}</p>
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
                    <div :class="['mcp-tab',{ 'active': tabActive === 0 }]" @click="tabClick(0)">介绍概览</div>
                    <div style="display: inline-block" v-if="hosted === '1'">
                        <div :class="['mcp-tab',{ 'active': tabActive === 1 }]" @click="tabClick(1)">SSE URL及工具</div>
                    </div>
                </div>

                <div v-if="tabActive === 0">
                    <div class="overview bg-border" >
                        <div class="overview-item">
                            <div class="item-title">• &nbsp;使用概述</div>
                            <div class="item-desc" v-html="parseTxt(detail.whatIs)"></div>
                        </div>
                        <div class="overview-item">
                            <div class="item-title">• &nbsp;特性说明</div>
                            <div class="item-desc" v-html="parseTxt(detail.keyFeatures)"></div>
                        </div>
                        <div class="overview-item">
                            <div class="item-title">• &nbsp;应用场景</div>
                            <div class="item-desc" >
                                <div v-html="parseTxt(detail.whereToUse)"></div>
                                <div v-html="parseTxt(detail.useCases)"></div>
                            </div>
                        </div>
                    </div>
                    <div class="overview bg-border" >
                        <div class="overview-item">
                            <div class="item-title">• &nbsp;使用说明</div>
                            <div class="item-desc" v-html="parseTxt(detail.howToUse)"></div>
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
                <!--<div class="install bg-border" v-if="tabActive === 2">
                    &lt;!&ndash;copy from https://mcpmarket.cn/&ndash;&gt;
                    <div class="login-required-message" style="text-align: center;background-color: #a7535305; padding: 40px 20px; border-radius: 8px; margin: 20px 0;">
                        <i class="fas fa-lock el-icon-lock" style="font-size: 48px; color: #D33A3A; margin-bottom: 20px; display: block;"></i>
                        <h3 style="margin-bottom: 15px; color: #333;font-size: 20px;">需要登录</h3>
                        <p style="margin-bottom: 25px; color: #666; line-height: 40px;">
                            要获取SSE URL和配置MCP服务器，请先登录您的账号。如果没有账号，您可以快速注册一个。
                        </p>
                        <div style="display: flex; justify-content: center; gap: 15px;">
                            <a href="https://mcpmarket.cn/auth/login?next=%2Fserver%2F67ff4974764487b6b9e11c21" style="display: inline-block; padding: 10px 20px; background-color: #D33A3A; color: white; text-decoration: none; border-radius: 6px; font-weight: 500;">
                                登录
                            </a>
                            <a href="https://mcpmarket.cn/auth/login?next=%2Fserver%2F67ff4974764487b6b9e11c21" style="display: inline-block; padding: 10px 20px; background-color: white; color: #D33A3A; text-decoration: none; border-radius: 6px; border: 1px solid #D33A3A; font-weight: 500;">
                                使用社交账号登录
                            </a>
                        </div>
                    </div>
                </div>-->

                <div class="tool bg-border" v-if="tabActive === 1">
                    <div class="tool-item ">
                        <p class="title">SSE URL:</p>
                        <div class="sse-url" style="display: flex">
                            <div class="tool-item-bg sse-url__input">{{detail.sseUrl}}</div>
                            <el-button class="sse-url__bt" type="primary" :disabled="!detail.isCanSendCustomMcp" @click="preSendToCustomize">发送到自定义</el-button>
                        </div>
                        <p style="line-height: 40px;color: #666;">* 将MCP发送到自定义后，您可在工作流或智能体中直接调用。</p>
                    </div>
                    <div class="tool-item" v-if="detail.tools">
                        <p class="title">工具介绍:</p>
                        <div class="tool-item-bg tool-intro">
                            <el-collapse class="mcp-el-collapse">
                                <el-collapse-item v-for="(n,i) in detail.tools" :title="n.toolName" :name="i">
                                    <div class="desc">描述：{{detail.description}}</div>
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

            <div class="right-recommend">
                <p style="margin: 20px 0;color: #333;">其他MCP服务查看</p>
                <div class="recommend-item" v-for="(n,i) in recommendList" :key="`${i}rc`" @click="handleClick(n)">
                    <img class="logo" :src="n.logo"/>
                    <p class="name">{{n.name}}</p>
                    <p class="intro">{{n.description}}</p>
                </div>
            </div>
        </div>

        <sendDialog ref="dialog" :dialogVisible="dialogVisible" :detail="detail" @handleClose="handleClose" @getIsCanSendStatus="getIsCanSendStatus"/>
    </div>
</template>
<script>
    import sendDialog from './sendDialog'
    import {md} from '../../mixins/marksown-it'
    import { getPublicMcpList, getMarkDownContent, getTools, getRecommendsList, getPublicMcpInfo } from "@/api/mcp";
    export default {
        data() {
            return {
                md:md,
                from: '',
                mcpSquareId:'',
                hosted:'',
                detail: {
                    description:'',
                    whatIs:'',
                    keyFeatures:'',
                    whereToUse:'',
                    useCases:'',
                    howToUse:'',
                    isCanSendCustomMcp:false
                },
                foldStatus:false,
                tabActive:0,
                recommendList: [
                  {
                    "mcpSquareId": "67e5dc0848048b1e353cb6af",
                    "name": "mcp-obsidian",
                    "description": "读取和搜索包含Markdown笔记的目录，专门为Obsidian vault设计。",
                    "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67e5dc0848048b1e353cb6af.png",
                    "hosted": false
                  },
                  {
                    "mcpSquareId": "67e5dc0848048b1e353cb6c1",
                    "name": "mcp-server-browserbase",
                    "description": "允许大型语言模型通过Browserbase和Stagehand控制浏览器",
                    "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67e5dc0848048b1e353cb6c1.png",
                    "hosted": false
                  },
                  {
                    "mcpSquareId": "67e5dc0a48048b1e353cb6eb",
                    "name": "mcp-jetbrains",
                    "description": "促进客户端与JetBrains IDE（如IntelliJ、PyCharm、WebStorm和Android Studio）之间的通信。",
                    "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67e5dc0a48048b1e353cb6eb.png",
                    "hosted": false
                  },
                  {
                    "mcpSquareId": "67e5dc0a48048b1e353cb6ed",
                    "name": "mcp-send-email",
                    "description": "通过Resend的API直接从Cursor发送电子邮件",
                    "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67e5dc0a48048b1e353cb6ed.png",
                    "hosted": false
                  },
                  {
                    "mcpSquareId": "67e5dc2448048b1e353cb99c",
                    "name": "playwright-mcp",
                    "description": "微软官方的Playwright工具",
                    "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67e5dc2448048b1e353cb99c.png",
                    "hosted": true
                  },
                  {
                    "mcpSquareId": "67e5dc2e48048b1e353cba91",
                    "name": "mcp-ical",
                    "description": "管理 macOS 日历, 包括安排会议、检查可用性、修改事件",
                    "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67e5dc2e48048b1e353cba91.png",
                    "hosted": false
                  },
                  {
                    "mcpSquareId": "67e5ddb748048b1e353cd81b",
                    "name": "paddle-mcp-server",
                    "description": "管理产品目录、账单和订阅，以及报告。",
                    "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67e5ddb748048b1e353cd81b.png",
                    "hosted": false
                  },
                  {
                    "mcpSquareId": "67f1294690965e7bf66c5e93",
                    "name": "brave-search",
                    "description": "Brave Search MCP Server集成了灵活过滤的网页和本地搜索功能。",
                    "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67f1294690965e7bf66c5e93.png",
                    "hosted": true
                  },
                  {
                    "mcpSquareId": "67f41180b66f446c3d8f8d7f",
                    "name": "elevenlabs-mcp",
                    "description": "官方 ElevenLabs MCP Server",
                    "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67f41180b66f446c3d8f8d7f.png",
                    "hosted": false
                  },
                  {
                    "mcpSquareId": "67f77bf1d6110df54f87a6e6",
                    "name": "MiniMax-MCP",
                    "description": " MiniMax官方MCP Server，支持高质量的视频生成、图像生成、语音生成、和声音克隆等多项能力",
                    "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67f77bf1d6110df54f87a6e6.png",
                    "hosted": true
                  }
                ],
                dialogVisible:false,
                markdownHtml:''
            };
        },
        watch: {
            $route: {
                handler: function(val, oldVal){
                    this.initData()
                },
                // 深度观察监听
                deep: true
            }
        },
        mounted() {
            this.initData()
            this.getRecommendList()
        },
        methods: {
            initData(){
                this.mcpSquareId = this.$route.params.id
                this.hosted = this.$route.params.hosted
                this.tabActive = 0
                this.getDataFromSquare()

                //滚动到顶部
                document.getElementById("timeScroll").scrollTop = 0
            },
            getDataFromSquare(){
                getPublicMcpInfo({mcpSquareId:this.mcpSquareId})
                    .then((res) => {
                        this.detail = res.data;
                        this.loadMarkDown()
                    })
                    .catch((err) => {
                      this.detail = {
                        "mcpSquareId": "67ff4974764487b6b9e11c21",
                        "alias": "",
                        "name": "Amap 高德地图",
                        "by": "高德开放平台",
                        "createdAt": "2025-06-11 16:08:35",
                        "description": "高德地图 MCP Server 现已覆盖12大核心服务接口，提供全场景覆盖的地图服务，包括地理编码、逆地理编码、IP 定位、天气查询、骑行路径规划、步行路径规划、驾车路径规划、公交路径规划、距离测量、关键词搜索、周边搜索、详情搜索等。",
                        "hosted": true,
                        "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67ff4974764487b6b9e11c21.png",
                        "modifiedBy": "admin",
                        "orderWeight": 6,
                        "port": 18000,
                        "sseUrl": "https://mcpmarket.cn/sse/67ff4974764487b6b9e11c21",
                        "howToUse": "\t•\t获取开发者 Key：登录高德开放平台并创建应用获取 Key\nhttps://lbs.amap.com/?ref=https://console.amap.com/dev/index\n\t•\t配置 MCP Server：在支持 MCP 的客户端（如 Cursor）中设置 SSE 或 Node.js 接入方式\n\t•\t连接模型：选择大模型（如 Claude），使用 Agent 模式进行交互\n\t•\t直接调用服务：通过快捷键打开交互窗口开始使用，如路线规划、美食推荐等",
                        "keyFeatures": "\t•\t零部署，易使用：无需本地服务器部署，仅通过配置 URL 即可使用。\n\t•\t语义优化结果：对返回 JSON 数据进行语义增强，便于大模型理解。\n\t•\t自动升级：平台持续迭代更新，无需用户手动操作。\n\t•\t全托管云服务：无需用户关注服务器维护或扩容。\n\t•\t协议兼容性强：支持标准 SSE 长连接协议，适配多种场景。",
                        "useCases": "\t•\t城市出行路线规划（骑行、步行、驾车、公交）\n\t•\t获取实时天气信息\n\t•\tIP 定位与地理编码/逆地理编码\n\t•\t地点搜索与 POI 信息查询\n\t•\t测量两点间距离",
                        "whatIs": "高德地图 MCP Server 是基于 SSE（Server-Sent Events）技术的地理服务接口集合，允许开发者通过 MCP 协议调用地图服务，如路径规划、天气查询、地点搜索等。它支持与如 Cursor、Claude 等大模型工具无缝集成。",
                        "whereToUse": "\t•\t在支持 MCP 协议的客户端中使用（如：Cursor、Claude、Cline）\n\t•\t可嵌入到企业应用、AI 助手系统、智能出行方案中\n\t•\t适用于城市交通服务平台、天气播报系统、位置服务 App 等",
                        "categories": [
                          "搜索",
                          "官方"
                        ],
                        "tags": [],
                        "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/67ff4974764487b6b9e11c21.md",
                        "tools": [
                          {
                            "toolName": "maps_regeocode",
                            "description": "描述：将一个高德经纬度坐标转换为行政区划地址信息",
                            "params": [
                              {
                                "name": "location",
                                "requiredBadge": "必填",
                                "type": "类型: string",
                                "description": "描述: 经纬度"
                              }
                            ]
                          },
                          {
                            "toolName": "maps_geo",
                            "description": "描述：将详细的结构化地址转换为经纬度坐标。支持对地标性名胜景区、建筑物名称解析为经纬度坐标",
                            "params": [
                              {
                                "name": "address",
                                "requiredBadge": "必填",
                                "type": "类型: string",
                                "description": "描述: 待解析的结构化地址信息"
                              },
                              {
                                "name": "city",
                                "type": "类型: string",
                                "description": "描述: 指定查询的城市"
                              }
                            ]
                          },
                          {
                            "toolName": "maps_ip_location",
                            "description": "描述：IP 定位根据用户输入的 IP 地址，定位 IP 的所在位置",
                            "params": [
                              {
                                "name": "ip",
                                "requiredBadge": "必填",
                                "type": "类型: string",
                                "description": "描述: IP地址"
                              }
                            ]
                          },
                          {
                            "toolName": "maps_weather",
                            "description": "描述：根据城市名称或者标准adcode查询指定城市的天气",
                            "params": [
                              {
                                "name": "city",
                                "requiredBadge": "必填",
                                "type": "类型: string",
                                "description": "描述: 城市名称或者adcode"
                              }
                            ]
                          },
                          {
                            "toolName": "maps_search_detail",
                            "description": "描述：查询关键词搜或者周边搜获取到的POI ID的详细信息",
                            "params": [
                              {
                                "name": "id",
                                "requiredBadge": "必填",
                                "type": "类型: string",
                                "description": "描述: 关键词搜或者周边搜获取到的POI ID"
                              }
                            ]
                          },
                          {
                            "toolName": "maps_bicycling",
                            "description": "描述：骑行路径规划用于规划骑行通勤方案，规划时会考虑天桥、单行线、封路等情况。最大支持 500km 的骑行路线规划",
                            "params": [
                              {
                                "name": "origin",
                                "requiredBadge": "必填",
                                "type": "类型: string",
                                "description": "描述: 出发点经纬度，坐标格式为：经度，纬度"
                              },
                              {
                                "name": "destination",
                                "requiredBadge": "必填",
                                "type": "类型: string",
                                "description": "描述: 目的地经纬度，坐标格式为：经度，纬度"
                              }
                            ]
                          },
                          {
                            "toolName": "maps_direction_walking",
                            "description": "描述：步行路径规划 API 可以根据输入起点终点经纬度坐标规划100km 以内的步行通勤方案，并且返回通勤方案的数据",
                            "params": [
                              {
                                "name": "origin",
                                "requiredBadge": "必填",
                                "type": "类型: string",
                                "description": "描述: 出发点经度，纬度，坐标格式为：经度，纬度"
                              },
                              {
                                "name": "destination",
                                "requiredBadge": "必填",
                                "type": "类型: string",
                                "description": "描述: 目的地经度，纬度，坐标格式为：经度，纬度"
                              }
                            ]
                          },
                          {
                            "toolName": "maps_direction_driving",
                            "description": "描述：驾车路径规划 API 可以根据用户起终点经纬度坐标规划以小客车、轿车通勤出行的方案，并且返回通勤方案的数据。",
                            "params": [
                              {
                                "name": "origin",
                                "requiredBadge": "必填",
                                "type": "类型: string",
                                "description": "描述: 出发点经度，纬度，坐标格式为：经度，纬度"
                              },
                              {
                                "name": "destination",
                                "requiredBadge": "必填",
                                "type": "类型: string",
                                "description": "描述: 目的地经度，纬度，坐标格式为：经度，纬度"
                              }
                            ]
                          },
                          {
                            "toolName": "maps_direction_transit_integrated",
                            "description": "描述：公交路径规划 API 可以根据用户起终点经纬度坐标规划综合各类公共（火车、公交、地铁）交通方式的通勤方案，并且返回通勤方案的数据，跨城场景下必须传起点城市与终点城市",
                            "params": [
                              {
                                "name": "origin",
                                "requiredBadge": "必填",
                                "type": "类型: string",
                                "description": "描述: 出发点经度，纬度，坐标格式为：经度，纬度"
                              },
                              {
                                "name": "destination",
                                "requiredBadge": "必填",
                                "type": "类型: string",
                                "description": "描述: 目的地经度，纬度，坐标格式为：经度，纬度"
                              },
                              {
                                "name": "city",
                                "requiredBadge": "必填",
                                "type": "类型: string",
                                "description": "描述: 公共交通规划起点城市"
                              },
                              {
                                "name": "cityd",
                                "requiredBadge": "必填",
                                "type": "类型: string",
                                "description": "描述: 公共交通规划终点城市"
                              }
                            ]
                          },
                          {
                            "toolName": "maps_distance",
                            "description": "描述：距离测量 API 可以测量两个经纬度坐标之间的距离,支持驾车、步行以及球面距离测量",
                            "params": [
                              {
                                "name": "origins",
                                "requiredBadge": "必填",
                                "type": "类型: string",
                                "description": "描述: 起点经度，纬度，可以传多个坐标，使用竖线隔离，比如120,30|120,31，坐标格式为：经度，纬度"
                              },
                              {
                                "name": "destination",
                                "requiredBadge": "必填",
                                "type": "类型: string",
                                "description": "描述: 终点经度，纬度，坐标格式为：经度，纬度"
                              },
                              {
                                "name": "type",
                                "type": "类型: string",
                                "description": "描述: 距离测量类型,1代表驾车距离测量，0代表直线距离测量，3步行距离测量"
                              }
                            ]
                          },
                          {
                            "toolName": "maps_text_search",
                            "description": "描述：关键词搜，根据用户传入关键词，搜索出相关的POI",
                            "params": [
                              {
                                "name": "keywords",
                                "requiredBadge": "必填",
                                "type": "类型: string",
                                "description": "描述: 搜索关键词"
                              },
                              {
                                "name": "city",
                                "type": "类型: string",
                                "description": "描述: 查询城市"
                              },
                              {
                                "name": "types",
                                "type": "类型: string",
                                "description": "描述: POI类型，比如加油站"
                              }
                            ]
                          },
                          {
                            "toolName": "maps_around_search",
                            "description": "描述：周边搜，根据用户传入关键词以及坐标location，搜索出radius半径范围的POI",
                            "params": [
                              {
                                "name": "keywords",
                                "type": "类型: string",
                                "description": "描述: 搜索关键词"
                              },
                              {
                                "name": "location",
                                "requiredBadge": "必填",
                                "type": "类型: string",
                                "description": "描述: 中心点经度纬度"
                              },
                              {
                                "name": "radius",
                                "type": "类型: string",
                                "description": "描述: 搜索半径"
                              }
                            ]
                          }
                        ],
                        "isCanSendCustomMcp": true
                      }
                    });
            },
            getIsCanSendStatus(){
                getPublicMcpInfo({mcpSquareId:this.mcpSquareId})
                    .then((res) => {
                        this.detail.isCanSendCustomMcp = res.data.isCanSendCustomMcp;
                    })
                    .catch((err) => {});
            },
            getRecommendList(){
                let params = {
                    mcpSquareId:this.mcpSquareId
                }
                getRecommendsList(params)
                    .then((res) => {
                        this.recommendList = res.data.list;
                    })
                    .catch((err) => {});
            },
            handleClick(val){
                this.$router.push(`/publicMCP/detail/${val.mcpSquareId}/${val.hosted?1:0}`)
            },
            //解析文本，遇到.换行等
            parseTxt(txt){
                return txt.replaceAll('\n\t','<br/>&nbsp;').replaceAll('\t', '   &nbsp;')
            },
            tabClick(status){
                this.tabActive = status
            },
            fold(){
                this.foldStatus = !this.foldStatus
            },
            preSendToCustomize(){
                this.dialogVisible = true
                this.$refs.dialog.ruleForm.serverUrl = this.detail.sseUrl
            },
            handleClose(){
                this.dialogVisible = false
            },
            back(){
                this.$router.push('/mcp')
            },
            //读取markdown文档
            loadMarkDown(){
                getMarkDownContent(this.detail.contentPath.split('https://obs-nmhhht6.cucloud.cn/')[1])
                    .then(async(res) => {
                        const text = await res.text();
                        this.markdownHtml = text
                    })
                    .catch((err) => {});
            },
        },
        components: {
            sendDialog
        },
    };
</script>
<style lang="scss">
@import "./markdown.min.css";
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
          flex:1;
          margin-right: 20px;
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
      cursor: pointer;
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


.mcp-el-collapse.el-collapse {
  border: none;
}
.mcp-el-collapse .el-collapse-item {
  margin: 10px 0;
  border: none;

  .el-collapse-item__header {
    border: none;
    color: $color;
    font-weight: bold;
    padding: 0 20px;
  }

  .el-collapse-item__wrap {
    padding: 0 20px;
    border: none;
  }

  .desc {
    background: #F4F5FF; // rgba(255, 246, 246, 1);
    padding: 10px 15px;
    border-radius: 6px;
    border: 1px solid #98A6E9; // #f5cbcb;
  }

  .params {
    margin-top: 12px;

    .params-table {
      border-radius: 6px;
      border: 1px solid #ddd;
      padding: 10px 12px;
      background-color: #fff;
      margin-top: 6px;

      .tr {
        display: flex;

        .td {
          padding: 0 30px 0 0;
        }

        .color {
          color: $color;
        }
      }

      .params-desc {
        margin-top: 4px;
        color: #999;
      }
    }
  }
}
.mcp-markdown {
  /deep/.code-header{
    height: 0!important;
    padding:0;
  }
}
.el-button.is-disabled {
  background: #f9f9f9 !important;
}
</style>