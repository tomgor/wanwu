<template>
  <div class="main-content docs-page-content">
    <div class="mark__content">
      <div class="markdown-body" v-html="mdContent"></div>
    </div>
  </div>
</template>
<script>

const highlight = require('highlight.js')
import 'highlight.js/styles/github.css'
import { marked } from 'marked'
import { getMarkdown } from "@/api/docs"
import { DOC_FIRST_KEY } from "../constants"
import { Loading } from 'element-ui'

marked.setOptions({
  renderer: new marked.Renderer(),
  gfm: true,
  tables: true,
  breaks: false,
  pedantic: false,
  sanitize: false,
  smartLists: true,
  smartypants: false,
  highlight: function (code) {
    return highlight.highlightAuto(code).value
  }
})

export default {
  data(){
    return {
      mdContent: '',
    }
  },
  watch: {
    $route: {
      handler (val, oldValue) {
        if (val !== oldValue) {
          this.getMarkdown(val.params.id)
          const docsPageContent = document.querySelector('.el-main')
          if (docsPageContent) docsPageContent.scrollTo(0, 0)
        }
      },
      // 深度观察监听
      deep: true
    }
  },
  created() {
    this.getMarkdown(this.$route.params.id)
  },
  methods: {
    getMarkdown(path) {
      if (path === DOC_FIRST_KEY) return

      const loadingInstance = Loading.service({
        target: document.querySelector('.docs-page-content')
      })
      getMarkdown({path}).then(res => {
        const mdContent = marked(res.data || this.$t('common.noData'))
        this.mdContent =  mdContent
        loadingInstance.close()
      }).catch(() => {
        const {data} = {
          "code": 0,
          "data": "# 元景大模型平台介绍\r\n\r\n中国联通元景大模型MaaS平台（以下简称“平台”），以选模型-\u003e改模型-\u003e用模型，构建联通特色的大模型工具链范式，一站式使能大模型开发应用，供行业大模型开发；同时提供基于智能体的原生应用开发工具，实现低门槛应用开发。\r\n\r\n平台为开发行业模型、行业应用提供选模型-\u003e改模型-\u003e用模型的全套联通特色工具链支持：\r\n\r\n1.在选模型方面，提供对比评测和自主评测工具，对大模型进行反馈效果对比，或批量数据集进行推理效果评测。用于用户进行基本模型的选择。\r\n\r\n2.在改模型方面，提供模型管理工具。支持用户基于基础大模型一键快速预训练、模型微调以及模型压缩。\r\n\r\n3.在用模型方面，提供0代码应用研发工具，使用应用研发工具快速创建在线应用进行使用和发布，同时支持知识库管理、提示词优化等。\r\n\r\n\r\n\r\n## **平台产品架构图及操作流程**\r\n\r\n**大模型MaaS平台产品架构图：**\r\n\r\n![](../../../service/api/v1/doc-center/v2.4.3-1747119948649/assets/image-20241203100756408.png)\r\n\r\n​                                                                                                    \r\n\r\n**平台首页介绍：**\r\n\r\n![](../../../service/api/v1/doc-center/v2.4.3-1747119948649/assets/image-20250114132636206-1736832399427-1.png)\r\n\r\n\r\n\r\n### 选模型模块核心功能\r\n\r\n**1）自主评测**\r\n\r\n用户可通过自主评测功能，在选模型阶段，对大模型进行效果预检测：通过指定格式文档上传可批量导入进行批量评测和评测结果标注。\r\n\r\n**2）模型比对**\r\n\r\n用户可同时选择多个模型，进行输入请求，查看多模型反馈效果，进行模型选用的评估参考依据。\r\n\r\n**3）在线体验** \r\n\r\n用户可在线体验各参数级大模型。体验界面同时可查看该模型的基本信息，便于用户使用\r\n\r\n### 改模型模块核心功能\r\n\r\n**1）模型预训练**\r\n\r\n平台提供基础大模型，用户可根据自身业务需要，选择不同的模型类别，导入数据集，进行模型预训练。支持断点续训。\r\n\r\n**2）SFT训练**\r\n\r\n用户自行创建微调任务，选择模型类别、微调形式、自建微调数据集、微调参数设置、进行模型训练。\r\n\r\n**3）DPO训练** 用户可对微调过后的模型进行DPO训练。可以帮助在保持模型性能的同时，减少资源消耗和提高部署效率。\r\n\r\n**4）RLHF训练** 用户可对微调过后的模型进行RLHF训练。\r\n\r\n**5）模型压缩**\r\n\r\n用户可对预训练和微调过后的模型进行量化压缩。可以帮助在保持模型性能的同时，减少资源消耗和提高部署效率。\r\n\r\n**6）训练任务管理**\r\n\r\n训练任务管理均在任务列表完成，可查看模型训练状态、训练进度。训练完毕的模型可进行模型发布。\r\n\r\n**7）模型管理**\r\n\r\n模型管理可查看所有已发布好的模型，版本描述及模型试用状态。训练完毕的模型可进行启动、在线测试、继续训练等操作。\r\n\r\n### 用模型模块核心功能\r\n\r\n**1）大模型在线使用**\r\n\r\n用户可在线体验大模型。平台提供几十种对话用提示词模版,如：邮件撰写、小红书模版、计划书模版等。用户可直接使用模版，替换关键词进行对话。\r\n\r\n**2）知识库管理**\r\n\r\n提供知识库管理，用户可进行多级知识库创建和4种类型文档上传，自定义设置文档分段形式。\r\n\r\n**3）原生应用开发工具**\r\n\r\n提供0代码原生应用开发工具，用户可使用推荐语、知识增强等开发能力，在线进行应用开发。开发测试完毕后可直接在线使用。 研发完毕的应用，同时可进行公开发布或私密发布，与他人共享使用。\r\n\r\n**4）原生应用超市**\r\n\r\n用户可在原生应用超市中选择需要使用的应用，原生应用超市中提供已创建好的应用，无需单独部署调用，可直接在线使用。\r\n\r\n**5）模型服务**\r\n\r\n用户在“工具箱-应用管理”中，创建好应用后，可获得应用ID，API KEY，密钥等信息，即可通过API进行模型调用，对训练好的模型进行部署。",
          "msg": ""
        }
        this.mdContent = marked(data) || this.$t('common.noData')
        loadingInstance.close()
      })
    }
  }
}
</script>

<style lang="scss" scoped>
@import "@/assets/showDocs/showdoc.scss";
.docs-page-content {
  margin: 0 0 50px;
}
</style>
