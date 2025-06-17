<template>
  <div class="mcp-content-box mcp-third">
    <div class="mcp-main">
      <div class="mcp-content">
        <!--<div class="mcp-menu">
          <p style="margin: 10px;border-bottom: 1px solid #d9d9d9;font-weight: bold">分类筛选</p>
          <el-radio-group v-model="category" @input="radioChange">
            <el-radio-button v-for="(n,i) in menuList" :key="`${i}ml`" :label="n"></el-radio-button>
          </el-radio-group>
        </div>-->
        <div class="mcp-card-box">
          <div class="card-search card-search-cust">
            <!--<div>
              按服务方式&nbsp;&nbsp;&nbsp;
              <el-radio-group v-model="hosted" @input="handleSearch" size="mini">
                <el-radio label="">不限</el-radio>
                <el-radio label="1">仅看云部署</el-radio>
                <el-radio label="0">仅看本地化</el-radio>
              </el-radio-group>
            </div>-->
            <div>
              <span
                v-for="item in typeList"
                :key="item.key"
                :class="['tab-span', {'is-active': typeRadio === item.key}]"
                @click="changeTab(item.key)"
              >
                {{item.name}}
              </span>
            </div>
            <search-input placeholder="请输入MCP名称进行搜索" ref="searchInput" @handleSearch="handleSearch" />
          </div>

          <div class="card-loading-box" v-if="list.length > 0">
            <div class="card-box"
              v-loading="listLoading"
              v-infinite-scroll="loadList"
              infinite-scroll-delay="200"
              :infinite-scroll-disabled="disabled"
            >
              <div
                class="card"
                v-for="(item, index) in list"
                :key="index"
                @click.stop="handleClick(item)"
              >
                <div class="card-title">
                  <img class="card-logo" :src="`${item.logo}?x-oss-process=image/resize,w_100`"/>
                  <div class="mcp_detailBox">
                    <span class="mcp_name">{{ item.alias || item.name }}</span>
                    <span class="mcp_from">
                      <label>
                        {{ item.by }}
                      </label>
                    </span>
                  </div>
                </div>
                <div class="card-des">{{ item.description }}</div>
                <!--<div :class="['hosted',item.hosted?'sse':'local']">{{item.hosted?'云部署':'本地化'}}</div>-->
              </div>
            </div>
            <p class="loading-tips" v-if="loading">加载中...</p>
            <!--<p class="loading-tips" v-if="noMore">没有更多了</p>-->
          </div>
          <div v-else class="empty"><el-empty description="暂无数据"></el-empty></div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import { getPublicMcpList, getMcpCategoryList } from "@/api/mcp"
import SearchInput from "@/components/searchInput.vue"
export default {
  components: { SearchInput },
  data() {
    return {
      mcpSquareId: "",
      name: "",
      hosted:'',
      category:'全部',
      menuList:[],
      list: [
        {
          "mcpSquareId": "682f31450b63a9cd4a7b0208",
          "alias": "",
          "name": "百度网盘",
          "by": "百度网盘",
          "createdAt": "2025-06-11 16:08:35",
          "description": "百度网盘MCP server涵盖用户信息、获取文件信息、上传下载、文件管理、文件搜索等。提供了SSE和stdio两种接入方式供用户选择。其中，上传功能仅可通过stdio进行本地接入，其余功能支持通过SSE方式接入。\n\n",
          "hosted": false,
          "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/682f31450b63a9cd4a7b0208.png",
          "modifiedBy": "admin",
          "orderWeight": 0,
          "port": 0,
          "sseUrl": "https://mcpmarket.cn/sse/682f31450b63a9cd4a7b0208",
          "howToUse": "使用百度网盘，用户首先需要注册百度账号并完成实名认证。在开放平台上创建应用后，用户可以获得服务密钥（AK）和访问令牌，以便与MCP服务器进行交互。",
          "keyFeatures": "百度网盘的关键特性包括获取文件列表、文档/图片/视频管理、创建文件夹、文件复制/移动/删除/重命名、本地文件上传、文件下载和文件搜索。它支持两种集成方式：SSE和stdio。",
          "useCases": "百度网盘的使用场景包括重要文件的个人备份、与同事分享大文件、团队项目文档的管理以及从不同设备远程访问文件。",
          "whatIs": "百度网盘是一个云存储服务，允许用户在线存储、管理和分享文件。它提供了一个MCP服务器，使客户端能够连接并执行各种操作，如文件检索、上传/下载和管理。",
          "whereToUse": "百度网盘广泛应用于个人文件存储、商业文档管理以及需要文件共享和访问的协作工作环境等多个领域。",
          "categories": [
            "最新",
            "数据",
            "搜索",
            "官方"
          ],
          "tags": [],
          "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/682f31450b63a9cd4a7b0208.md",
          "tools": null
        },
        {
          "mcpSquareId": "680f7f400b63a9cd4a538b25",
          "alias": "Skywork Mureka MCP",
          "name": "Mureka-mcp",
          "by": "SkyworkAI",
          "createdAt": "2025-06-11 16:08:35",
          "description": "生成歌词、歌曲和背景音乐（伴奏）。模型上下文协议（Model Context Protocol，MCP）服务器。",
          "hosted": false,
          "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/680f7f400b63a9cd4a538b25.png",
          "modifiedBy": "admin",
          "orderWeight": 0,
          "port": 0,
          "sseUrl": "https://mcpmarket.cn/sse/680f7f400b63a9cd4a538b25",
          "howToUse": "使用 Mureka-mcp，首先从 Mureka 平台获取 API 密钥，安装 'uv' 包管理器，并在您的 MCP 客户端（如 Claude Desktop）中配置必要的设置，包括您的 API 密钥和 URL。重启客户端以访问可用的 MCP 工具。",
          "keyFeatures": "Mureka-mcp 的关键特性包括生成歌词、歌曲和器乐背景音乐的能力，与各种 MCP 客户端的无缝集成，以及用户友好的设置过程。",
          "useCases": "Mureka-mcp 的使用场景包括为艺术家创作原创歌曲，为视频和游戏生成背景音乐，以及帮助词曲作者进行歌词和旋律的头脑风暴。",
          "whatIs": "Mureka-mcp 是一个官方的模型上下文协议（MCP）服务器，能够通过强大的 API 生成歌词、歌曲和背景音乐（器乐）。它允许与各种 MCP 客户端交互，创建音乐内容。",
          "whereToUse": "Mureka-mcp 可用于音乐制作、内容创作、游戏开发以及任何需要自动生成音乐和歌词的应用领域。",
          "categories": [
            "社交",
            "创作",
            "官方"
          ],
          "tags": [],
          "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/680f7f400b63a9cd4a538b25.md",
          "tools": null
        },
        {
          "mcpSquareId": "68016281a5deaa7aaa926ac4",
          "alias": "",
          "name": "azure-mcp",
          "by": "Azure",
          "createdAt": "2025-06-11 16:08:35",
          "description": "This repository is for development of the Azure MCP Server, bringing the power of Azure to your agents.",
          "hosted": false,
          "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/68016281a5deaa7aaa926ac4.png",
          "modifiedBy": "",
          "orderWeight": 0,
          "port": 0,
          "sseUrl": "https://mcpmarket.cn/sse/68016281a5deaa7aaa926ac4",
          "howToUse": "要使用 Azure MCP，您可以通过在 Visual Studio Code 中执行以下命令来安装它：``npx -y @azure/mcp@latest server start``。",
          "keyFeatures": "关键特性包括探索 Azure 资源、查询和分析数据、管理配置以及执行高级 Azure 操作。",
          "useCases": "使用场景包括查询 Azure 存储帐户、管理 Cosmos DB 数据库、分析日志数据，以及使用 Node.js 构建 Azure 应用程序。",
          "whatIs": "Azure MCP 是一个旨在通过与 Azure 服务的无缝集成来增强 AI 代理的服务器，遵循 MCP 规范。",
          "whereToUse": "Azure MCP 可用于云计算、数据分析、应用开发和 AI 集成等多个领域。",
          "categories": [
            "社交",
            "生产",
            "数据"
          ],
          "tags": [],
          "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/68016281a5deaa7aaa926ac4.md",
          "tools": null
        },
        {
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
          ]
        },
        {
          "mcpSquareId": "67ff0119764487b6b9e0da90",
          "alias": "Alipay 支付宝",
          "name": "Alipay MCP",
          "by": "Alipay",
          "createdAt": "2025-06-11 16:08:35",
          "description": "支付宝官方MCP Server，让你可以轻松将支付宝开放平台提供的交易创建、查询、退款等能力集成到你的 LLM 应用中",
          "hosted": false,
          "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67ff0119764487b6b9e0da90.png",
          "modifiedBy": "admin",
          "orderWeight": 23,
          "port": 0,
          "sseUrl": "https://mcpmarket.cn/sse/67ff0119764487b6b9e0da90",
          "howToUse": "要使用 Alipay MCP，你需要先成为支付宝开放平台的商户并获取商户私钥。然后，在你的项目设置中配置 MCP Server，并通过主流的 MCP Client 使用它。",
          "keyFeatures": "Alipay MCP 的关键特性包括与支付宝支付服务的无缝集成、自动生成支付链接，以及通过 AI 驱动的交互提升交易效率。",
          "useCases": "一个典型的使用场景是插画师利用 Alipay MCP 创建智能聊天应用，客户可以描述插画需求，快速获得报价并支付，无需人工干预，从而简化整个交易流程。",
          "whatIs": "Alipay MCP 是支付宝开放平台提供的一个服务器，允许将交易创建、查询、退款等能力轻松集成到你的 LLM 应用中，从而开发具备支付功能的智能工具。",
          "whereToUse": "Alipay MCP 可应用于电子商务、创意服务等多个领域，适用于任何需要支付处理和交易管理的应用。",
          "categories": [
            "金融",
            "官方",
            "电商"
          ],
          "tags": [],
          "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/67ff0119764487b6b9e0da90.md",
          "tools": null
        },
        {
          "mcpSquareId": "67fb5b712e1080dbe767acbf",
          "alias": "Skyvern",
          "name": "skyvern",
          "by": "Skyvern-AI",
          "createdAt": "2025-06-11 16:08:35",
          "description": "    使用LLMs和计算机视觉自动化基于浏览器的工作流程",
          "hosted": false,
          "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67fb5b712e1080dbe767acbf.png",
          "modifiedBy": "admin",
          "orderWeight": 15,
          "port": 0,
          "sseUrl": "https://mcpmarket.cn/sse/67fb5b712e1080dbe767acbf",
          "howToUse": "使用 Skyvern，开发者可以将其 API 集成到他们的应用程序中，从而在网站上自动化任务，无需复杂的脚本或维护。",
          "keyFeatures": "Skyvern 的关键特性包括手动工作流程的自动化、用户友好的 API，以及能够处理各种网站而不依赖脆弱脚本的能力。",
          "useCases": "Skyvern 的使用场景包括自动化表单提交、从网页提取数据、在电子商务网站上执行重复任务，以及与其他软件工具集成以提高生产力。",
          "whatIs": "Skyvern 是一个利用大型语言模型（LLMs）和计算机视觉来自动化基于浏览器的工作流程的平台，提供简单的 API 来简化各种网站上的手动流程。",
          "whereToUse": "Skyvern 可以应用于多个领域，例如网页抓取、数据录入自动化、测试网页应用程序，以及任何需要浏览器交互的场景。",
          "categories": [
            "开发",
            "官方"
          ],
          "tags": [
            "api",
            "automation",
            "browser",
            "computer",
            "gpt",
            "llm",
            "playwright",
            "python",
            "rpa",
            "vision",
            "workflow",
            "browser-automation"
          ],
          "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/67fb5b712e1080dbe767acbf.md",
          "tools": null
        },
        {
          "mcpSquareId": "67fb5a202e1080dbe767ab53",
          "alias": "MarkItDown  ",
          "name": "markitdown",
          "by": "microsoft",
          "createdAt": "2025-06-11 16:08:35",
          "description": "将文件和办公文档转换为Markdown的Python工具。",
          "hosted": false,
          "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67fb5a202e1080dbe767ab53.png",
          "modifiedBy": "admin",
          "orderWeight": 14,
          "port": 0,
          "sseUrl": "https://mcpmarket.cn/sse/67fb5a202e1080dbe767ab53",
          "howToUse": "要使用 MarkItDown，可以通过命令 'pip install markitdown' 安装。然后，您可以使用 DocumentConverter 类将类文件对象转换为 Markdown 格式，确保使用二进制类文件对象以保持兼容性。",
          "keyFeatures": "MarkItDown 的关键特性包括在转换为 Markdown 时能够保留重要的文档结构，如标题、列表、表格和链接。它还支持可选功能组以增强功能。",
          "useCases": "MarkItDown 的使用场景包括将办公文档转换为 LLM 应用中的分析格式，为机器学习任务准备文档，以及促进从各种文件格式中提取结构化数据。",
          "whatIs": "MarkItDown 是一个轻量级的 Python 工具，旨在将各种文件和办公文档转换为 Markdown 格式，主要用于 LLM 应用和文本分析管道。",
          "whereToUse": "MarkItDown 非常适合用于需要文本分析的领域，如数据科学、自然语言处理以及涉及 LLM 和文档转换的任何应用。",
          "categories": [
            "工作",
            "创作",
            "官方"
          ],
          "tags": [
            "langchain",
            "openai",
            "autogen-extension",
            "autogen",
            "markdown",
            "microsoft-office",
            "pdf"
          ],
          "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/67fb5a202e1080dbe767ab53.md",
          "tools": null
        },
        {
          "mcpSquareId": "67fb54bb2e1080dbe767a4ff",
          "alias": "",
          "name": "HyperChat",
          "by": "BigSweetPotatoStudio",
          "createdAt": "2025-06-11 16:08:35",
          "description": "HyperChat 是一个聊天客户端，致力于开放性，利用各种 LLMs（大型语言模型）的 API 来实现最佳的聊天体验，同时通过 MCP 协议实现生产力工具。",
          "hosted": false,
          "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67fb54bb2e1080dbe767a4ff.png",
          "modifiedBy": "",
          "orderWeight": 0,
          "port": 0,
          "sseUrl": "https://mcpmarket.cn/sse/67fb54bb2e1080dbe767a4ff",
          "howToUse": "使用HyperChat，您可以通过命令行运行，使用命令``npx -y @dadigua/hyper-chat``，通过浏览器访问``http://localhost:16100/123456/``，或使用Docker，命令为``docker pull dadigua/hyperchat-mini:latest``。",
          "keyFeatures": "HyperChat的关键特性包括支持多种操作系统（Windows、MacOS、Linux）、命令行操作、Docker兼容性、WebDAV同步、MCP扩展、黑暗模式、资源和提示支持、多会话聊天、定时任务以及聊天中的模型比较。",
          "useCases": "HyperChat的使用场景包括团队协作、自动化任务管理、教育辅导和实时客户服务互动。",
          "whatIs": "HyperChat是一个开源聊天客户端，旨在实现开放性，利用各种LLM的API提供最佳聊天体验，并通过MCP协议实现生产力工具。",
          "whereToUse": "HyperChat可以用于软件开发、客户支持、教育等多个领域，任何需要高效沟通和生产力工具的地方。",
          "categories": [
            "社交",
            "生产",
            "创作"
          ],
          "tags": [
            "chat-application",
            "llm",
            "mcp",
            "modelcontextprotocol",
            "agent",
            "local-agent"
          ],
          "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/67fb54bb2e1080dbe767a4ff.md",
          "tools": null
        },
        {
          "mcpSquareId": "67fb4da22e1080dbe7679b53",
          "alias": "",
          "name": "inbox-zero",
          "by": "elie222",
          "createdAt": "2025-06-11 16:08:35",
          "description": "AI个人助手用于电子邮件。开源应用程序帮助您快速达到收件箱清空的目标。",
          "hosted": false,
          "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67fb4da22e1080dbe7679b53.png",
          "modifiedBy": "",
          "orderWeight": 0,
          "port": 0,
          "sseUrl": "https://mcpmarket.cn/sse/67fb4da22e1080dbe7679b53",
          "howToUse": "用户可以通过将 Inbox Zero 与他们的电子邮件账户集成来使用它。AI助手根据用户定义的提示进行操作，可以执行诸如草拟回复、归档邮件和管理垃圾邮件等操作。",
          "keyFeatures": "关键特性包括用于邮件管理的AI个人助手、用于跟踪需要回复的邮件的Reply Zero、联系人智能分类、快速管理邮件的批量退订功能、冷邮件拦截器以及用于跟踪活动的邮件分析。",
          "useCases": "使用场景包括管理大量电子邮件、自动化重复的电子邮件任务、跟踪重要对话、退订不需要的新闻通讯以及分析电子邮件参与度指标。",
          "whatIs": "Inbox Zero 是一个AI个人助手，旨在帮助用户高效管理电子邮件，并快速实现“收件箱清零”的状态。它是一个开源应用程序，可以自动化各种电子邮件任务。",
          "whereToUse": "Inbox Zero 可以在多个领域使用，特别是在电子邮件沟通频繁的专业环境中，如企业环境、客户服务和个人生产力。",
          "categories": [
            "生产",
            "社交",
            "创作"
          ],
          "tags": [
            "ai",
            "email",
            "nextjs",
            "openai",
            "prisma",
            "productivity",
            "tailwind",
            "typescript",
            "turborepo",
            "shadcn-ui",
            "upstash",
            "gmail",
            "loops",
            "postgresql",
            "posthog",
            "resend",
            "tinybird"
          ],
          "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/67fb4da22e1080dbe7679b53.md",
          "tools": null
        },
        {
          "mcpSquareId": "67fa946b2e1080dbe7670350",
          "alias": "微信读书MCP",
          "name": "mcp-server-weread",
          "by": "freestylefly",
          "createdAt": "2025-06-11 16:08:35",
          "description": "一个桥接微信读书数据和Claude Desktop的轻量级服务器，使您可以在Claude中无缝访问微信读书的笔记和阅读数据。",
          "hosted": false,
          "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67fa946b2e1080dbe7670350.png",
          "modifiedBy": "admin",
          "orderWeight": 0,
          "port": 0,
          "sseUrl": "https://mcpmarket.cn/sse/67fa946b2e1080dbe7670350",
          "howToUse": "使用mcp-server-weread时，请确保已安装Node.js（v16+），克隆仓库，安装依赖，获取微信读书Cookie，配置环境变量并启动服务器。最后，配置您的MCP客户端以连接到服务器。",
          "keyFeatures": "关键特性包括获取用户书架、获取带有笔记的书籍、获取特定书籍的笔记、获取书籍详细信息、按关键词搜索笔记以及获取最近阅读的书籍。",
          "useCases": "使用场景包括查看书架、查看特定书籍的笔记、搜索包含特定关键词的笔记以及获取最近阅读的书籍。",
          "whatIs": "mcp-server-weread是一个轻量级服务器，旨在桥接微信读书数据和Claude Desktop，使用户能够在Claude中无缝访问阅读笔记和数据。",
          "whereToUse": "mcp-server-weread可用于教育技术、个人知识管理和AI辅助阅读应用，特别适合使用微信读书和Claude Desktop的用户。",
          "categories": [
            "社交",
            "创作"
          ],
          "tags": [],
          "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/67fa946b2e1080dbe7670350.md",
          "tools": null
        },
        {
          "mcpSquareId": "67f7e358cbe9764bc98f6484",
          "alias": "RedNote 小红书",
          "name": "RedNote-MCP",
          "by": "iFurySt",
          "createdAt": "2025-06-11 16:08:35",
          "description": "用于访问小红书，可用关键词搜索笔记，访问帖子内容等。",
          "hosted": false,
          "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67f7e358cbe9764bc98f6484.png",
          "modifiedBy": "admin",
          "orderWeight": 27,
          "port": 0,
          "sseUrl": "https://mcpmarket.cn/sse/67f7e358cbe9764bc98f6484",
          "howToUse": "使用 RedNote-MCP，首先通过 npm 全局安装，命令为 'npm install -g rednote-mcp'。然后运行 'rednote-mcp init' 初始化登录，这将打开浏览器让你登录小红书并保存会话 Cookie。",
          "keyFeatures": "关键特性包括支持 Cookie 持久化的认证管理、关键词搜索笔记、通过 URL 访问笔记内容、通过 URL 访问评论内容，以及命令行初始化工具。",
          "useCases": "RedNote-MCP 的使用场景包括构建需要小红书内容的应用、对用户生成内容进行研究，以及开发内容管理和分析工具。",
          "whatIs": "RedNote-MCP 是一个用于访问小红书（XiaoHongShu, xhs）内容的 Model Context Protocol 服务器。它通过结构化协议促进与平台数据的交互。",
          "whereToUse": "RedNote-MCP 可用于数据分析、内容聚合和应用开发等多个领域，特别是在需要访问小红书内容的场景中。",
          "categories": [
            "社交",
            "搜索"
          ],
          "tags": [
            "ai",
            "mcp",
            "mcp-server",
            "rednote",
            "xhs",
            "rednote-mcp",
            "xhs-mcp"
          ],
          "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/67f7e358cbe9764bc98f6484.md",
          "tools": null
        },
        {
          "mcpSquareId": "67f77bf1d6110df54f87a6e6",
          "alias": "MiniMax",
          "name": "MiniMax-MCP",
          "by": "MiniMax-AI",
          "createdAt": "2025-06-11 16:08:35",
          "description": " MiniMax官方MCP Server，支持高质量的视频生成、图像生成、语音生成、和声音克隆等多项能力",
          "hosted": true,
          "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67f77bf1d6110df54f87a6e6.png",
          "modifiedBy": "admin",
          "orderWeight": 4,
          "port": 18007,
          "sseUrl": "https://mcpmarket.cn/sse/67f77bf1d6110df54f87a6e6",
          "howToUse": "使用MiniMax-MCP，开发者可以将提供的API集成到他们的应用程序中，从而实现无缝的文本转语音转换和视频生成功能。",
          "keyFeatures": "MiniMax-MCP的关键特性包括高质量的文本转语音合成、高效的视频生成能力，以及便于集成的用户友好API。",
          "useCases": "MiniMax-MCP的使用场景包括创建教育视频、为多媒体项目生成配音、开发互动学习工具，以及提升数字内容的无障碍性。",
          "whatIs": "MiniMax-MCP是MiniMax模型上下文协议的官方服务器，旨在促进与先进的文本转语音和视频生成API的交互。",
          "whereToUse": "MiniMax-MCP可以应用于教育、娱乐、内容创作和为残疾人士提供无障碍解决方案等多个领域。",
          "categories": [
            "创作",
            "官方"
          ],
          "tags": [],
          "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/67f77bf1d6110df54f87a6e6.md",
          "tools": [
            {
              "toolName": "text_to_audio",
              "description": "描述：Convert text to audio with a given voice and save the output audio file to a given directory.\n    Directory is optional, if not provided, the output file will be saved to $HOME/Desktop.\n    Voice id is optional, if not provided, the default voice will be used.\n\n    COST WARNING: This tool makes an API call to Minimax which may incur costs. Only use when explicitly requested by the user.\n\n    Args:\n        text (str): The text to convert to speech.\n        voice_id (str, optional): The id of the voice to use. For example, \"male-qn-qingse\"/\"audiobook_female_1\"/\"cute_boy\"/\"Charming_Lady\"...\n        model (string, optional): The model to use.\n        speed (float, optional): Speed of the generated audio. Controls the speed of the generated speech. Values range from 0.5 to 2.0, with 1.0 being the default speed. \n        vol (float, optional): Volume of the generated audio. Controls the volume of the generated speech. Values range from 0 to 10, with 1 being the default volume.\n        pitch (int, optional): Pitch of the generated audio. Controls the speed of the generated speech. Values range from -12 to 12, with 0 being the default speed.\n        emotion (str, optional): Emotion of the generated audio. Controls the emotion of the generated speech. Values range [\"happy\", \"sad\", \"angry\", \"fearful\", \"disgusted\", \"surprised\", \"neutral\"], with \"happy\" being the default emotion.\n        sample_rate (int, optional): Sample rate of the generated audio. Controls the sample rate of the generated speech. Values range [8000,16000,22050,24000,32000,44100] with 32000 being the default sample rate.\n        bitrate (int, optional): Bitrate of the generated audio. Controls the bitrate of the generated speech. Values range [32000,64000,128000,256000] with 128000 being the default bitrate.\n        channel (int, optional): Channel of the generated audio. Controls the channel of the generated speech. Values range [1, 2] with 1 being the default channel.\n        format (str, optional): Format of the generated audio. Controls the format of the generated speech. Values range [\"pcm\", \"mp3\",\"flac\"] with \"mp3\" being the default format.\n        language_boost (str, optional): Language boost of the generated audio. Controls the language boost of the generated speech. Values range ['Chinese', 'Chinese,Yue', 'English', 'Arabic', 'Russian', 'Spanish', 'French', 'Portuguese', 'German', 'Turkish', 'Dutch', 'Ukrainian', 'Vietnamese', 'Indonesian', 'Japanese', 'Italian', 'Korean', 'Thai', 'Polish', 'Romanian', 'Greek', 'Czech', 'Finnish', 'Hindi', 'auto'] with \"auto\" being the default language boost.\n        output_directory (str): The directory to save the audio to.\n\n    Returns:\n        Text content with the path to the output file and name of the voice used.\n    ",
              "params": [
                {
                  "name": "text",
                  "requiredBadge": "必填",
                  "type": "类型: string",
                  "description": "描述: "
                },
                {
                  "name": "output_directory",
                  "type": "类型: string",
                  "description": "描述: "
                },
                {
                  "name": "voice_id",
                  "type": "类型: string",
                  "default": "默认值: female-shaonv",
                  "description": "描述: "
                },
                {
                  "name": "model",
                  "type": "类型: string",
                  "default": "默认值: speech-02-hd",
                  "description": "描述: "
                },
                {
                  "name": "speed",
                  "type": "类型: number",
                  "default": "默认值: 1.0",
                  "description": "描述: "
                },
                {
                  "name": "vol",
                  "type": "类型: number",
                  "default": "默认值: 1.0",
                  "description": "描述: "
                },
                {
                  "name": "pitch",
                  "type": "类型: integer",
                  "description": "描述: "
                },
                {
                  "name": "emotion",
                  "type": "类型: string",
                  "default": "默认值: happy",
                  "description": "描述: "
                },
                {
                  "name": "sample_rate",
                  "type": "类型: integer",
                  "default": "默认值: 32000",
                  "description": "描述: "
                },
                {
                  "name": "bitrate",
                  "type": "类型: integer",
                  "default": "默认值: 128000",
                  "description": "描述: "
                },
                {
                  "name": "channel",
                  "type": "类型: integer",
                  "default": "默认值: 1",
                  "description": "描述: "
                },
                {
                  "name": "format",
                  "type": "类型: string",
                  "default": "默认值: mp3",
                  "description": "描述: "
                },
                {
                  "name": "language_boost",
                  "type": "类型: string",
                  "default": "默认值: auto",
                  "description": "描述: "
                }
              ]
            },
            {
              "toolName": "list_voices",
              "description": "描述：List all voices available.\n\n    Args:\n        voice_type (str, optional): The type of voices to list. Values range [\"all\", \"system\", \"voice_cloning\"], with \"all\" being the default.\n    Returns:\n        Text content with the list of voices.\n    ",
              "params": [
                {
                  "name": "voice_type",
                  "type": "类型: string",
                  "default": "默认值: all",
                  "description": "描述: "
                }
              ]
            },
            {
              "toolName": "voice_clone",
              "description": "描述：Clone a voice using provided audio files. The new voice will be charged upon first use.\n\n    COST WARNING: This tool makes an API call to Minimax which may incur costs. Only use when explicitly requested by the user.\n\n     Args:\n        voice_id (str): The id of the voice to use.\n        file (str): The path to the audio file to clone or a URL to the audio file.\n        text (str, optional): The text to use for the demo audio.\n        is_url (bool, optional): Whether the file is a URL. Defaults to False.\n        output_directory (str): The directory to save the demo audio to.\n    Returns:\n        Text content with the voice id of the cloned voice.\n    ",
              "params": [
                {
                  "name": "voice_id",
                  "requiredBadge": "必填",
                  "type": "类型: string",
                  "description": "描述: "
                },
                {
                  "name": "file",
                  "requiredBadge": "必填",
                  "type": "类型: string",
                  "description": "描述: "
                },
                {
                  "name": "text",
                  "requiredBadge": "必填",
                  "type": "类型: string",
                  "description": "描述: "
                },
                {
                  "name": "output_directory",
                  "type": "类型: string",
                  "description": "描述: "
                },
                {
                  "name": "is_url",
                  "type": "类型: boolean",
                  "description": "描述: "
                }
              ]
            },
            {
              "toolName": "play_audio",
              "description": "描述：Play an audio file. Supports WAV and MP3 formats. Not supports video.\n\n     Args:\n        input_file_path (str): The path to the audio file to play.\n        is_url (bool, optional): Whether the audio file is a URL.\n    Returns:\n        Text content with the path to the audio file.\n    ",
              "params": [
                {
                  "name": "input_file_path",
                  "requiredBadge": "必填",
                  "type": "类型: string",
                  "description": "描述: "
                },
                {
                  "name": "is_url",
                  "type": "类型: boolean",
                  "description": "描述: "
                }
              ]
            },
            {
              "toolName": "generate_video",
              "description": "描述：Generate a video from a prompt.\n\n    COST WARNING: This tool makes an API call to Minimax which may incur costs. Only use when explicitly requested by the user.\n\n     Args:\n        model (str, optional): The model to use. Values range [\"T2V-01\", \"T2V-01-Director\", \"I2V-01\", \"I2V-01-Director\", \"I2V-01-live\"]. \"Director\" supports inserting instructions for camera movement control. \"I2V\" for image to video. \"T2V\" for text to video.\n        prompt (str): The prompt to generate the video from. When use Director model, the prompt supports 15 Camera Movement Instructions (Enumerated Values)\n            -Truck: [Truck left], [Truck right]\n            -Pan: [Pan left], [Pan right]\n            -Push: [Push in], [Pull out]\n            -Pedestal: [Pedestal up], [Pedestal down]\n            -Tilt: [Tilt up], [Tilt down]\n            -Zoom: [Zoom in], [Zoom out]\n            -Shake: [Shake]\n            -Follow: [Tracking shot]\n            -Static: [Static shot]\n        first_frame_image (str): The first frame image. The model must be \"I2V\" Series.\n        output_directory (str): The directory to save the video to.\n        async_mode (bool, optional): Whether to use async mode. Defaults to False. If True, the video generation task will be submitted asynchronously and the response will return a task_id. Should use ``query_video_generation`` tool to check the status of the task and get the result.\n    Returns:\n        Text content with the path to the output video file.\n    ",
              "params": [
                {
                  "name": "model",
                  "type": "类型: string",
                  "default": "默认值: T2V-01",
                  "description": "描述: "
                },
                {
                  "name": "prompt",
                  "type": "类型: string",
                  "description": "描述: "
                },
                {
                  "name": "first_frame_image",
                  "type": "类型: string",
                  "description": "描述: "
                },
                {
                  "name": "output_directory",
                  "type": "类型: string",
                  "description": "描述: "
                },
                {
                  "name": "async_mode",
                  "type": "类型: boolean",
                  "description": "描述: "
                }
              ]
            },
            {
              "toolName": "query_video_generation",
              "description": "描述：Query the status of a video generation task.\n\n    Args:\n        task_id (str): The task ID to query. Should be the task_id returned by ``generate_video`` tool if ``async_mode`` is True.\n        output_directory (str): The directory to save the video to.\n    Returns:\n        Text content with the status of the task.\n    ",
              "params": [
                {
                  "name": "task_id",
                  "requiredBadge": "必填",
                  "type": "类型: string",
                  "description": "描述: "
                },
                {
                  "name": "output_directory",
                  "type": "类型: string",
                  "description": "描述: "
                }
              ]
            },
            {
              "toolName": "text_to_image",
              "description": "描述：Generate a image from a prompt.\n\n    COST WARNING: This tool makes an API call to Minimax which may incur costs. Only use when explicitly requested by the user.\n\n     Args:\n        model (str, optional): The model to use. Values range [\"image-01\"], with \"image-01\" being the default.\n        prompt (str): The prompt to generate the image from.\n        aspect_ratio (str, optional): The aspect ratio of the image. Values range [\"1:1\", \"16:9\",\"4:3\", \"3:2\", \"2:3\", \"3:4\", \"9:16\", \"21:9\"], with \"1:1\" being the default.\n        n (int, optional): The number of images to generate. Values range [1, 9], with 1 being the default.\n        prompt_optimizer (bool, optional): Whether to optimize the prompt. Values range [True, False], with True being the default.\n        output_directory (str): The directory to save the image to.\n    Returns:\n        Text content with the path to the output image file.\n    ",
              "params": [
                {
                  "name": "model",
                  "type": "类型: string",
                  "default": "默认值: image-01",
                  "description": "描述: "
                },
                {
                  "name": "prompt",
                  "type": "类型: string",
                  "description": "描述: "
                },
                {
                  "name": "aspect_ratio",
                  "type": "类型: string",
                  "default": "默认值: 1:1",
                  "description": "描述: "
                },
                {
                  "name": "n",
                  "type": "类型: integer",
                  "default": "默认值: 1",
                  "description": "描述: "
                },
                {
                  "name": "prompt_optimizer",
                  "type": "类型: boolean",
                  "default": "默认值: True",
                  "description": "描述: "
                },
                {
                  "name": "output_directory",
                  "type": "类型: string",
                  "description": "描述: "
                }
              ]
            }
          ]
        },
        {
          "mcpSquareId": "67f6b924d6110df54f86ddc4",
          "alias": "SearchAPI",
          "name": "searchAPI-mcp",
          "by": "RmMargt",
          "createdAt": "2025-06-11 16:08:35",
          "description": "作为AI助手与搜索服务之间的桥梁，支持地图搜索、航班查询、酒店预订等功能。",
          "hosted": false,
          "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67f6b924d6110df54f86ddc4.png",
          "modifiedBy": "admin",
          "orderWeight": 45,
          "port": 0,
          "sseUrl": "https://mcpmarket.cn/sse/67f6b924d6110df54f86ddc4",
          "howToUse": "使用 searchAPI-mcp，开发者可以通过向服务器发出 API 调用将其集成到应用程序中。他们可以通过统一的接口访问地图搜索、航班查询和酒店预订等多种搜索功能，从而轻松实现搜索能力。",
          "keyFeatures": "searchAPI-mcp 的关键特性包括网页搜索结果、知识图谱集成、多语言支持、视频内容搜索、Google Maps 地点搜索、航班搜索选项和酒店可用性查询。它还提供时间范围、视频时长等过滤选项。",
          "useCases": "searchAPI-mcp 的使用场景包括需要航班和酒店搜索的旅行预订应用程序、提供基于位置服务的虚拟助手，以及集成视频内容搜索功能的教育平台。",
          "whatIs": "searchAPI-mcp 是一个基于 Model Context Protocol (MCP) 的搜索 API 服务器，提供对 Google Maps、Google Flights 和 Google Hotels 等多种 Google 服务的标准化访问。它充当 AI 助手与搜索服务之间的桥梁，实现无缝交互。",
          "whereToUse": "searchAPI-mcp 可用于旅行规划、酒店业、电子商务等多个领域，适用于任何需要地图、航班和酒店搜索功能的应用程序。它特别适合需要为用户提供相关搜索结果的 AI 助手和聊天机器人。",
          "categories": [
            "开发",
            "搜索"
          ],
          "tags": [],
          "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/67f6b924d6110df54f86ddc4.md",
          "tools": null
        },
        {
          "mcpSquareId": "67f4a521d6110df54f8411ca",
          "alias": "ChatPPT",
          "name": "chatppt-mcp",
          "by": "YOOTeam",
          "createdAt": "2025-06-11 16:08:35",
          "description": "# 基于 ChatPPT 的 AI PPT 生成服务\n\n该服务基于 ChatPPT，用于生成 AI 驱动的 PPT。它可以根据主题或需求创建演示文稿，也可以从上传的文档中生成。它支持在线编辑和下载。",
          "hosted": false,
          "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67f4a521d6110df54f8411ca.png",
          "modifiedBy": "admin",
          "orderWeight": 51,
          "port": 0,
          "sseUrl": "https://mcpmarket.cn/sse/67f4a521d6110df54f8411ca",
          "howToUse": "使用 chatppt-mcp，用户可以输入所需的主题或上传文档。服务将生成一个 PPT，用户可以在线编辑并下载以供离线使用。",
          "keyFeatures": "chatppt-mcp 的关键特性包括 AI 驱动的 PPT 生成、支持文档上传、在线编辑功能和可下载的演示文稿。它还逐步开放了 STDIO 模式，并将很快支持 SSE 协议。",
          "useCases": "chatppt-mcp 的使用场景包括创建教育演示文稿、商业提案、市场营销演示以及基于上传文档或特定主题的报告。",
          "whatIs": "chatppt-mcp 是一个基于 chatppt 的 AI PPT 生成服务。用户可以根据主题或要求生成演示文稿，也可以通过上传文档进行 PPT 生成。该服务支持在线编辑和下载。",
          "whereToUse": "chatppt-mcp 可以用于教育、商业演示、市场营销等多个领域，适用于任何需要有效视觉沟通的场景。",
          "categories": [
            "社交",
            "工作",
            "创作",
            "官方"
          ],
          "tags": [
            "agent",
            "aippt",
            "chatppt",
            "mcp",
            "office"
          ],
          "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/67f4a521d6110df54f8411ca.md",
          "tools": null
        },
        {
          "mcpSquareId": "67f41180b66f446c3d8f8d7f",
          "alias": "ElevenLabs ",
          "name": "elevenlabs-mcp",
          "by": "elevenlabs",
          "createdAt": "2025-06-11 16:08:35",
          "description": "官方 ElevenLabs MCP Server",
          "hosted": false,
          "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67f41180b66f446c3d8f8d7f.png",
          "modifiedBy": "admin",
          "orderWeight": 61,
          "port": 0,
          "sseUrl": "https://mcpmarket.cn/sse/67f41180b66f446c3d8f8d7f",
          "howToUse": "使用elevenlabs-mcp，首先从ElevenLabs获取API密钥，安装'uv'包管理器，并配置您的客户端应用程序（如Claude Desktop）以使用提供的API密钥连接到MCP服务器。",
          "keyFeatures": "关键特性包括语音克隆、语音生成、音频转录，以及与Claude Desktop、Cursor、Windsurf和OpenAI Agents等多种MCP客户端的兼容性。",
          "useCases": "elevenlabs-mcp的使用场景包括为视频创建配音、生成个性化音频消息、转录会议记录，以及开发互动语音应用程序。",
          "whatIs": "elevenlabs-mcp是由ElevenLabs开发的官方模型上下文协议（MCP）服务器，旨在促进与先进的文本转语音和音频处理API的交互。",
          "whereToUse": "elevenlabs-mcp可用于软件开发、内容创作、无障碍解决方案以及任何需要音频处理或语音合成的应用领域。",
          "categories": [
            "创作",
            "官方"
          ],
          "tags": [
            "elevenlabs",
            "elevenlabs-api",
            "mcp"
          ],
          "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/67f41180b66f446c3d8f8d7f.md",
          "tools": null
        },
        {
          "mcpSquareId": "67f3f4a1b66f446c3d8f7408",
          "alias": "",
          "name": "witsy",
          "by": "nbonamy",
          "createdAt": "2025-06-11 16:08:35",
          "description": "# Witsy: 桌面 AI 助手",
          "hosted": false,
          "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67f3f4a1b66f446c3d8f7408.png",
          "modifiedBy": "",
          "orderWeight": 0,
          "port": 0,
          "sseUrl": "https://mcpmarket.cn/sse/67f3f4a1b66f446c3d8f7408",
          "howToUse": "用户可以从 witsyai.com 或发布页面下载 Witsy，设置所需 LLM 提供商的 API 密钥，或通过 Ollama 使用本地模型。",
          "keyFeatures": "Witsy 支持多种 AI 模型，包括 OpenAI、Ollama、Anthropic 等。它提供了聊天完成与视觉模型支持、文本转图像和文本转视频功能，以及图像处理等特性。",
          "useCases": "使用场景包括生成创意内容、自动化客户支持、增强教育工具和开发互动应用程序。",
          "whatIs": "Witsy 是一款桌面 AI 助手，允许用户通过提供自己的 API 密钥或通过 Ollama 使用本地模型来利用各种 AI 模型。",
          "whereToUse": "Witsy 可用于内容创作、软件开发、教育等多个领域，以及任何需要 AI 驱动助手的领域。",
          "categories": [
            "创作",
            "生产"
          ],
          "tags": [
            "anthropic",
            "genai",
            "groq",
            "ollama",
            "ollama-gui",
            "openai",
            "electron-app",
            "electronjs",
            "vuejs",
            "vuejs3",
            "deepseek",
            "gemini",
            "googleai",
            "mistral",
            "mistralai",
            "openrouter",
            "mcp",
            "mcp-client"
          ],
          "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/67f3f4a1b66f446c3d8f7408.md",
          "tools": null
        },
        {
          "mcpSquareId": "67f3ee2db66f446c3d8f6e80",
          "alias": "Weather",
          "name": "weather-mcp-server",
          "by": "TuanKiri",
          "createdAt": "2025-06-11 16:08:35",
          "description": "一个轻量级的模型上下文协议（MCP）服务器，使得像Claude这样的人工智能助手能够获取和解读实时天气数据。",
          "hosted": false,
          "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67f3ee2db66f446c3d8f6e80.png",
          "modifiedBy": "admin",
          "orderWeight": 30,
          "port": 0,
          "sseUrl": "https://mcpmarket.cn/sse/67f3ee2db66f446c3d8f6e80",
          "howToUse": "要在 Claude Desktop 中使用 weather-mcp-server，请通过指定命令路径和您的 WEATHER_API_KEY 将其添加到 Claude 配置中。如果需要，您也可以使用 Go 从源代码构建服务器。",
          "keyFeatures": "关键特性包括能够获取指定城市的当前天气数据，使 AI 助手能够提供最新的天气信息。",
          "useCases": "使用场景包括在虚拟助手应用中提供天气更新、将天气数据集成到客户服务聊天机器人中，以及增强旅行和规划应用中的用户体验。",
          "whatIs": "weather-mcp-server 是一个轻量级的模型上下文协议（MCP）服务器，旨在允许 AI 助手（如 Claude）访问和解释实时天气数据。",
          "whereToUse": "weather-mcp-server 可用于需要实时天气信息的应用程序，例如虚拟助手、聊天机器人以及任何将天气数据集成到其服务中的软件。",
          "categories": [
            "数据",
            "社交",
            "精选"
          ],
          "tags": [
            "go",
            "golang",
            "mcp",
            "mcp-server"
          ],
          "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/67f3ee2db66f446c3d8f6e80.md",
          "tools": null
        },
        {
          "mcpSquareId": "67f3a63cb66f446c3d8f0296",
          "alias": "Youtube",
          "name": "mcp-youtube",
          "by": "anaisbetts",
          "createdAt": "2025-06-11 16:08:35",
          "description": "# YouTube视频摘要和字幕提取总结",
          "hosted": false,
          "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67f3a63cb66f446c3d8f0296.png",
          "modifiedBy": "admin",
          "orderWeight": 34,
          "port": 0,
          "sseUrl": "https://mcpmarket.cn/sse/67f3a63cb66f446c3d8f0296",
          "howToUse": "使用 mcp-youtube，首先通过 Homebrew 或 WinGet 在本地安装 'yt-dlp'，然后使用 mcp-installer 安装服务器，名称为 '@anaisbetts/mcp-youtube'。",
          "keyFeatures": "关键特性包括使用 'yt-dlp' 下载 YouTube 视频的字幕，以及与 AI 模型进行上下文交互的集成。",
          "useCases": "使用场景包括通过询问 Claude.ai 来总结 YouTube 视频、提取字幕以便翻译，以及提高视频内容的可访问性。",
          "whatIs": "mcp-youtube 是一个为 YouTube 设计的 Model-Context Protocol Server，允许用户下载字幕并与像 Claude.ai 这样的 AI 模型进行交互。",
          "whereToUse": "mcp-youtube 可用于教育、内容创作和 AI 研究等多个领域，特别是在需要视频内容分析的场景中。",
          "categories": [
            "社交",
            "创作"
          ],
          "tags": [],
          "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/67f3a63cb66f446c3d8f0296.md",
          "tools": null
        },
        {
          "mcpSquareId": "67f3a518b66f446c3d8f012a",
          "alias": "",
          "name": "ClaudeComputerCommander",
          "by": "wonderwhy-er",
          "createdAt": "2025-06-11 16:08:35",
          "description": "这是一个MCP Server，提供终端控制、文件系统搜索和差异文件编辑功能，供Claude使用。",
          "hosted": false,
          "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67f3a518b66f446c3d8f012a.png",
          "modifiedBy": "",
          "orderWeight": 0,
          "port": 0,
          "sseUrl": "https://mcpmarket.cn/sse/67f3a518b66f446c3d8f012a",
          "howToUse": "用户可以安装该服务器，并通过 Claude 桌面应用程序直接执行终端命令、管理文件和进行编辑任务。",
          "keyFeatures": "关键特性包括支持输出流的终端命令执行、命令超时支持、后台执行、进程管理（列出和终止进程）以及长时间运行命令的会话管理。",
          "useCases": "使用场景包括自动化重复的终端任务、管理系统进程、编辑配置文件以及高效地执行文件搜索和更新。",
          "whatIs": "ClaudeComputerCommander 是一个为 Claude 桌面应用程序设计的 MCP 服务器，提供终端控制、文件系统搜索和差异文件编辑功能。",
          "whereToUse": "ClaudeComputerCommander 可用于软件开发、系统管理以及需要终端命令执行和文件管理的任何环境。",
          "categories": [
            "社交",
            "生产",
            "开发"
          ],
          "tags": [
            "agent",
            "ai",
            "code-analysis",
            "code-generation",
            "mcp",
            "terminal-ai",
            "terminal-automation",
            "vibe-coding"
          ],
          "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/67f3a518b66f446c3d8f012a.md",
          "tools": null
        },
        {
          "mcpSquareId": "67f39d6cb66f446c3d8ef88c",
          "alias": "GitLab ",
          "name": "gitlab",
          "by": "modelcontextprotocol",
          "createdAt": "2025-06-11 16:08:35",
          "description": "MCP Server 用于 GitLab API，支持项目管理、文件操作等功能。",
          "hosted": false,
          "logo": "https://obs-nmhhht6.cucloud.cn/maas-public/mcp_square/logo/67f39d6cb66f446c3d8ef88c.png",
          "modifiedBy": "admin",
          "orderWeight": 81,
          "port": 0,
          "sseUrl": "https://mcpmarket.cn/sse/67f39d6cb66f446c3d8ef88c",
          "howToUse": "使用GitLab时，您可以通过MCP Server创建或更新文件、一次性推送多个文件、搜索代码库以及创建新项目。每个操作需要特定的输入，如项目ID、文件路径和提交信息。",
          "keyFeatures": "关键特性包括自动创建分支、全面的错误处理、保持Git历史记录以及支持批量操作，使用户能够高效管理文件。",
          "useCases": "GitLab的使用场景包括软件项目的版本控制、协作编码、持续集成与部署以及项目文档的管理。",
          "whatIs": "GitLab是一个基于Web的DevOps生命周期工具，提供Git代码库管理功能，包括问题跟踪、持续集成/持续交付（CI/CD）和项目管理。MCP Server为GitLab API提供项目管理和文件操作的支持。",
          "whereToUse": "GitLab广泛应用于软件开发、项目管理和各行业的协作，包括技术、教育和研究等领域。",
          "categories": [
            "开发",
            "工作",
            "企业",
            "官方"
          ],
          "tags": [],
          "contentPath": "https://obs-nmhhht6.cucloud.cn/maas-public/resource/mcp_square/content_md/67f39d6cb66f446c3d8ef88c.md",
          "tools": null
        }
      ],
      page:{
        pageNo: 1,
        pageSize: 20,
        total:0
      },
      count:0,
      loading:false,
      listLoading:false,
      typeRadio: '',
      typeList: [
        {name: '全部', key: ''},
        {name: '数据', key: 'data'},
        {name: '创作', key: 'creator'},
        {name: '搜索', key: 'search'},
      ]
    };
  },
  mounted() {
    this.doGetPublicMcpList()
    // this.doGetMcpCategoryList()
  },
  computed: {
    noMore () {
      return this.count >= this.page.total
    },
    disabled () {
      return this.loading || this.noMore
    }
  },
  methods: {
    changeTab(key) {
      this.typeRadio = key
      this.handleSearch()
    },
    doGetMcpCategoryList(){
      this.listLoading = true
      getMcpCategoryList()
        .then((res) => {
          this.menuList = res.data.list;
          this.listLoading = false
        })
        .catch((err) => {});
    },
    loadList(){
      this.loading = true
      setTimeout(()=>{
        this.doGetPublicMcpList()
      },400)
    },
    doGetPublicMcpList(){
      const searchInput = this.$refs.searchInput
      let params = {
        name: searchInput.value,
        pageNo: this.page.pageNo,
        pageSize: this.page.pageSize,
        category:this.category,
        type: this.typeRadio
      }
      if(this.hosted!==''){
        params.hosted = this.hosted === '1' ? true : false
      }

      getPublicMcpList(params)
        .then((res) => {
          //this.list = res.data.list;

          //懒加载
          this.list = this.list.concat(res.data.list);
          this.count += res.data.list.length
          this.page.total = res.data.total

          this.loading = false
          if(this.count < res.data.total){
            this.page.pageNo ++;
          }
        })
        .catch((err) => {});
    },
    radioChange(val){
      this.initPage()
      this.name = "";
      this.doGetPublicMcpList()
    },
    handleSearch() {
      this.initPage()
      this.doGetPublicMcpList()
    },
    initPage(){
      this.page = {
        pageNo: 1,
        pageSize: 20,
        total:0
      };
      this.count = 0;
      this.loading = false;
      this.list = [];
    },
    handleClick(val) {
      this.mcpSquareId = val.mcpSquareId;
      this.$router.push({path:`/publicMCP/detail/${val.mcpSquareId}/${val.hosted?1:0}`})
    },
  },
};
</script>
<style lang="scss">
.mcp-management .mcp-third{
  min-height: 600px;
  .el-radio-button__inner{
    border: none!important;
  }
  .tab-span {
    display: inline-block;
    vertical-align: middle;
    padding: 6px 12px;
    border-radius: 6px;
    color: $color_title;
    cursor: pointer;
  }
  .tab-span.is-active {
    color: $color;
    background: #fff;
    font-weight: bold;
  }
  .mcp-main{
    display: flex;
    padding: 0 20px;
    height: 100%;
    .mcp-content{
      display: flex;
      width:100%;
      padding: 0;
      /* height: calc(100vh - 230px);
       max-height: 100%;*/
      height: 100%;
      .mcp-menu{
        margin-top: 10px;
        margin-right: 20px;
        width: 90px;
        height: 450px;
        border: 1px solid $border_color; //#d0a7a7
        text-align: center;
        border-radius: 6px;
        color: #333;
        p{
          line-height: 28px;
          margin:10px 0;
        }
        .active{
          background: rgba(253, 231, 231, 1);
        }
      }
      .mcp-card-box{
        /*width: calc(100% - 110px);*/
        width: 100%;
        height: 100%;
        .input-with-select {
          width: 300px;
        }
        .card-loading-box{
         /* height: calc(100% - 200px);
          overflow: auto;*/
          .card-box {
            /*height: calc(100% - 34px);*/
            .hosted{
              position: absolute;
              right: -8px;
              top:5px;
              padding: 2px 6px;
              font-size: 12px;
              border-radius: 2px;
            }
            .sse{
              background: #d81e06;
              color: #fff;
            }
            .local{
              background: #555;
              color: #fff;
            }
          }
          .loading-tips{
            height: 30px;
            color: #999;
            text-align: center;
          }
        }
      }
    }
  }

  .card-logo{
    width: 50px;
    height: 50px;
    object-fit: cover;
  }
}
.el-radio-button:first-child .el-radio-button__inner,
.el-radio-button:last-child .el-radio-button__inner{
  border-radius: 0!important;
}
</style>