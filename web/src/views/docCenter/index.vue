<template>
  <div class="doc-page-container">
    <div class="doc-header">
      <div class="doc-header__left">
        <img v-if="homeLogoPath" style="height: 50px;" :src="basePath + '/user/api' + homeLogoPath"/>
        <span v-if="homeTitle">
          {{homeTitle}}
        </span>
      </div>
      <div class="doc-header__right">
        <el-select
          filterable
          remote
          clearable
          placeholder=""
          :remote-method="visibleChange"
          style="width: 500px;"
          class="top-search-input"
          popper-class="top-search-popover"
          :popper-append-to-body="false"
          v-model="searchText"
        >
          <i
            slot="prefix"
            class="el-input__icon el-icon-search"
            style="font-weight: bolder; font-size: 16px;"
          />
          <div class="header_search-option" v-if="searchText">
            <el-option v-if="searchList" :value="searchText" style="height: auto; background: #fff">
              <div v-if="docLoading" style="text-align: center; padding: 50px 0">
                <i class="el-icon-loading" style="font-size: 28px; color: #df0000"></i>
              </div>
              <div v-if="!docLoading && !(searchList && searchList.length)" style="text-align: center; padding: 50px 0">
                <span style="font-size: 14px; font-weight: normal; color: #999">{{$t('header.noData')}}</span>
              </div>
              <div
                v-if="!docLoading && searchList && searchList.length"
                v-for="(item, index) in searchList"
                :key="`search${item.title + index}`"
              >
                <div class="header_search-title">{{item.title}}</div>
                <div class="header_search-item" v-for="(it, i) in item.list" :key="`it${it.title + i}`">
                  <div class="header_search-item-left">
                    {{it.title}}
                  </div>
                  <div class="header_search-item-right" @click="jumpMenu(it.url)" v-html="it.content" />
                </div>
              </div>
            </el-option>
          </div>
        </el-select>
      </div>
    </div>
    <div class="doc-outer-container">
      <div :class="['doc-inner-container']">
        <el-aside style="min-width: 200px; width: auto; max-width: 300px" class="full-menu-aside">
          <el-menu
            :default-openeds="defaultOpeneds"
            :default-active="activeIndex"
            class="el-menu-vertical-demo"
          >
            <div v-for="(n,i) in menuList" :key="`${i}ml`">
              <!--有下一级-->
              <el-submenu v-if="n.children && checkPerm(n.perm)" :index="n.index">
                <template slot="title">
                  <i :class="n.icon || 'el-icon-menu'"></i>
                  <span>{{n.name}}</span>
                </template>
                <div v-for="(m,j) in n.children" v-if="checkPerm(m.perm)" :key="`${j}cl`">
                  <el-submenu v-if="m.children" :index="m.index" class="menu-indent">
                    <template slot="title">{{m.name}}</template>
                    <div v-for="(p,k) in m.children" :key="`${k}pl`" v-if="checkPerm(p.perm)">
                      <el-submenu v-if="p.children" :index="p.index" class="menu-indent-sub">
                        <template slot="title">{{p.name}}</template>
                        <el-menu-item
                          v-for="(item, index) in p.children"
                          :key="`${index}itemEl`"
                          :index="item.index"
                          v-if="checkPerm(item.perm)"
                          @click="menuClick(item)"
                        >
                          {{item.name}}
                        </el-menu-item>
                      </el-submenu>
                      <el-menu-item v-else :index="p.index" @click="menuClick(p)">{{p.name}}</el-menu-item>
                    </div >
                  </el-submenu>
                  <el-menu-item v-else :index="m.index" @click="menuClick(m)" class="menu-indent-item">{{m.name}}</el-menu-item>
                </div >
              </el-submenu>
              <!--没有下一级-->
              <el-menu-item :index="n.index" v-if="!n.children && checkPerm(n.perm)" @click="menuClick(n)">
                <i :class="n.icon || 'el-icon-menu'"></i>
                <span slot="title">{{n.name}}</span>
              </el-menu-item>
            </div>
          </el-menu>
        </el-aside>
        <!-- 右侧内容 -->
        <div class="doc-page-main">
          <DocPage />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import DocPage from "./components/docPage.vue"
import { checkPerm } from "@/router/permission"
import { DOC_FIRST_KEY } from "./constants"
import { getDocMenu, getDocSearchContent } from "@/api/docs"
import { fetchPermFirPath, fetchCurrentPathIndex } from "@/utils/util"
import { mapGetters } from "vuex"

export default {
  components: { DocPage },
  data() {
    return{
      basePath: this.$basePath,
      homeLogoPath: '',
      homeTitle: '',
      defaultOpeneds: [],
      menuList: [],
      docMenuList: [],
      docLoading: false,
      searchList: [],
      searchText: '',
      activeIndex: '0',
    }
  },
  watch: {
    $route: {
      handler (val) {
        this.changeMenuIndex(fetchCurrentPathIndex(val.path, this.menuList))
      },
      // 深度观察监听
      deep: true
    },
    commonInfo:{
      handler(val) {
        const { home = {} } = val.data || {}
        this.homeLogoPath = home.logo ? home.logo.path : ''
        this.homeTitle = home.title || ''
      },
      deep: true
    }
  },
  computed: {
    ...mapGetters('user', ['commonInfo']),
  },
  async created() {
    // 获取菜单
    this.getCurrentMenu()
  },
  methods: {
    checkPerm,
    jumpMenu(url) {
      // location.href = url
      const [_, path] = url.split(`${this.$basePath}/aibase`)
      this.$router.push({path})
    },
    loadSearchMenu() {
      this.docLoading = true
      this.searchList = []
      getDocSearchContent({content: this.searchText}).then(res => {
        this.docLoading = false
        this.searchList = res.data || []
      }).catch(() => {
        this.docLoading = false
      })
    },
    visibleChange(val) {
      this.searchText = val
      if (val) {
        this.loadSearchMenu()
      }
    },
    menuClick(item) {
      if (item.redirect) {
        item.redirect()
      } else {
        // 文档中心返回不带页面 path 前缀，跳转加上 path 前缀，避免点击路径直接拼到当前链接后面等问题
        this.$router.push({path: `/docCenter/pages/${item.path}`})
      }
    },
    getDocMenu() {
      return getDocMenu().then(res => {
        this.docMenuList = res.data || []
      }).catch(() => {
        const {data} = {
          "code": 0,
          "data": [
            {
              "name": "1 元景大模型平台介绍",
              "index": "doc1",
              "path": "v2.4.3-1747119948649%2F1%20%E5%85%83%E6%99%AF%E5%A4%A7%E6%A8%A1%E5%9E%8B%E5%B9%B3%E5%8F%B0%E4%BB%8B%E7%BB%8D.md",
              "pathRaw": "v2.4.3-1747119948649/1 元景大模型平台介绍.md",
              "children": null
            },
            {
              "name": "2 平台操作说明",
              "index": "doc2",
              "path": "",
              "pathRaw": "",
              "children": [
                {
                  "name": "1 选模型",
                  "index": "doc2-1",
                  "path": "",
                  "pathRaw": "",
                  "children": [
                    {
                      "name": "模型广场",
                      "index": "doc2-1-1",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F1%20%E9%80%89%E6%A8%A1%E5%9E%8B%2F%E6%A8%A1%E5%9E%8B%E5%B9%BF%E5%9C%BA.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/1 选模型/模型广场.md",
                      "children": null
                    },
                    {
                      "name": "模型排行",
                      "index": "doc2-1-2",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F1%20%E9%80%89%E6%A8%A1%E5%9E%8B%2F%E6%A8%A1%E5%9E%8B%E6%8E%92%E8%A1%8C.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/1 选模型/模型排行.md",
                      "children": null
                    },
                    {
                      "name": "模型推荐",
                      "index": "doc2-1-3",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F1%20%E9%80%89%E6%A8%A1%E5%9E%8B%2F%E6%A8%A1%E5%9E%8B%E6%8E%A8%E8%8D%90.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/1 选模型/模型推荐.md",
                      "children": null
                    }
                  ]
                },
                {
                  "name": "2 改模型",
                  "index": "doc2-2",
                  "path": "",
                  "pathRaw": "",
                  "children": [
                    {
                      "name": "1 模型仓库",
                      "index": "doc2-2-1",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F2%20%E6%94%B9%E6%A8%A1%E5%9E%8B%2F1%20%E6%A8%A1%E5%9E%8B%E4%BB%93%E5%BA%93.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/2 改模型/1 模型仓库.md",
                      "children": null
                    },
                    {
                      "name": "2 模型训练",
                      "index": "doc2-2-2",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F2%20%E6%94%B9%E6%A8%A1%E5%9E%8B%2F2%20%E6%A8%A1%E5%9E%8B%E8%AE%AD%E7%BB%83.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/2 改模型/2 模型训练.md",
                      "children": null
                    },
                    {
                      "name": "3 模型压缩",
                      "index": "doc2-2-3",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F2%20%E6%94%B9%E6%A8%A1%E5%9E%8B%2F3%20%E6%A8%A1%E5%9E%8B%E5%8E%8B%E7%BC%A9.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/2 改模型/3 模型压缩.md",
                      "children": null
                    },
                    {
                      "name": "4 模型评估",
                      "index": "doc2-2-4",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F2%20%E6%94%B9%E6%A8%A1%E5%9E%8B%2F4%20%E6%A8%A1%E5%9E%8B%E8%AF%84%E4%BC%B0.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/2 改模型/4 模型评估.md",
                      "children": null
                    }
                  ]
                },
                {
                  "name": "3 用模型",
                  "index": "doc2-3",
                  "path": "",
                  "pathRaw": "",
                  "children": [
                    {
                      "name": "1 大模型在线使用",
                      "index": "doc2-3-1",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F3%20%E7%94%A8%E6%A8%A1%E5%9E%8B%2F1%20%E5%A4%A7%E6%A8%A1%E5%9E%8B%E5%9C%A8%E7%BA%BF%E4%BD%BF%E7%94%A8.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/3 用模型/1 大模型在线使用.md",
                      "children": null
                    },
                    {
                      "name": "2 原生应用开发平台",
                      "index": "doc2-3-2",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F3%20%E7%94%A8%E6%A8%A1%E5%9E%8B%2F2%20%E5%8E%9F%E7%94%9F%E5%BA%94%E7%94%A8%E5%BC%80%E5%8F%91%E5%B9%B3%E5%8F%B0.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/3 用模型/2 原生应用开发平台.md",
                      "children": null
                    },
                    {
                      "name": "3 提示词工程",
                      "index": "doc2-3-3",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F3%20%E7%94%A8%E6%A8%A1%E5%9E%8B%2F3%20%E6%8F%90%E7%A4%BA%E8%AF%8D%E5%B7%A5%E7%A8%8B.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/3 用模型/3 提示词工程.md",
                      "children": null
                    },
                    {
                      "name": "4 知识库管理",
                      "index": "doc2-3-4",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F3%20%E7%94%A8%E6%A8%A1%E5%9E%8B%2F4%20%E7%9F%A5%E8%AF%86%E5%BA%93%E7%AE%A1%E7%90%86.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/3 用模型/4 知识库管理.md",
                      "children": null
                    },
                    {
                      "name": "5 问题管理",
                      "index": "doc2-3-5",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F3%20%E7%94%A8%E6%A8%A1%E5%9E%8B%2F5%20%E9%97%AE%E9%A2%98%E7%AE%A1%E7%90%86.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/3 用模型/5 问题管理.md",
                      "children": null
                    }
                  ]
                },
                {
                  "name": "4 数据集管理",
                  "index": "doc2-4",
                  "path": "",
                  "pathRaw": "",
                  "children": [
                    {
                      "name": "1 训练数据集管理",
                      "index": "doc2-4-1",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F4%20%E6%95%B0%E6%8D%AE%E9%9B%86%E7%AE%A1%E7%90%86%2F1%20%E8%AE%AD%E7%BB%83%E6%95%B0%E6%8D%AE%E9%9B%86%E7%AE%A1%E7%90%86.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/4 数据集管理/1 训练数据集管理.md",
                      "children": null
                    },
                    {
                      "name": "2 数据处理",
                      "index": "doc2-4-2",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F4%20%E6%95%B0%E6%8D%AE%E9%9B%86%E7%AE%A1%E7%90%86%2F2%20%E6%95%B0%E6%8D%AE%E5%A4%84%E7%90%86.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/4 数据集管理/2 数据处理.md",
                      "children": null
                    },
                    {
                      "name": "3 数据标注",
                      "index": "doc2-4-3",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F4%20%E6%95%B0%E6%8D%AE%E9%9B%86%E7%AE%A1%E7%90%86%2F3%20%E6%95%B0%E6%8D%AE%E6%A0%87%E6%B3%A8.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/4 数据集管理/3 数据标注.md",
                      "children": null
                    }
                  ]
                },
                {
                  "name": "5 应用管理",
                  "index": "doc2-5",
                  "path": "",
                  "pathRaw": "",
                  "children": [
                    {
                      "name": "应用管理",
                      "index": "doc2-5-1",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F5%20%E5%BA%94%E7%94%A8%E7%AE%A1%E7%90%86%2F%E5%BA%94%E7%94%A8%E7%AE%A1%E7%90%86.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/5 应用管理/应用管理.md",
                      "children": null
                    }
                  ]
                },
                {
                  "name": "6 高级资源管理",
                  "index": "doc2-6",
                  "path": "",
                  "pathRaw": "",
                  "children": [
                    {
                      "name": "1 GPU资源管理",
                      "index": "doc2-6-1",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F6%20%E9%AB%98%E7%BA%A7%E8%B5%84%E6%BA%90%E7%AE%A1%E7%90%86%2F1%20GPU%E8%B5%84%E6%BA%90%E7%AE%A1%E7%90%86.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/6 高级资源管理/1 GPU资源管理.md",
                      "children": null
                    },
                    {
                      "name": "2 服务管理",
                      "index": "doc2-6-2",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F6%20%E9%AB%98%E7%BA%A7%E8%B5%84%E6%BA%90%E7%AE%A1%E7%90%86%2F2%20%E6%9C%8D%E5%8A%A1%E7%AE%A1%E7%90%86.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/6 高级资源管理/2 服务管理.md",
                      "children": null
                    },
                    {
                      "name": "3 推理管理",
                      "index": "doc2-6-3",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F6%20%E9%AB%98%E7%BA%A7%E8%B5%84%E6%BA%90%E7%AE%A1%E7%90%86%2F3%20%E6%8E%A8%E7%90%86%E7%AE%A1%E7%90%86.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/6 高级资源管理/3 推理管理.md",
                      "children": null
                    }
                  ]
                },
                {
                  "name": "7 权限管理",
                  "index": "doc2-7",
                  "path": "",
                  "pathRaw": "",
                  "children": [
                    {
                      "name": "1 角色管理",
                      "index": "doc2-7-1",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F7%20%E6%9D%83%E9%99%90%E7%AE%A1%E7%90%86%2F1%20%E8%A7%92%E8%89%B2%E7%AE%A1%E7%90%86.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/7 权限管理/1 角色管理.md",
                      "children": null
                    },
                    {
                      "name": "2 用户管理",
                      "index": "doc2-7-2",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F7%20%E6%9D%83%E9%99%90%E7%AE%A1%E7%90%86%2F2%20%E7%94%A8%E6%88%B7%E7%AE%A1%E7%90%86.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/7 权限管理/2 用户管理.md",
                      "children": null
                    },
                    {
                      "name": "3 组织管理",
                      "index": "doc2-7-3",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F7%20%E6%9D%83%E9%99%90%E7%AE%A1%E7%90%86%2F3%20%E7%BB%84%E7%BB%87%E7%AE%A1%E7%90%86.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/7 权限管理/3 组织管理.md",
                      "children": null
                    }
                  ]
                },
                {
                  "name": "8 控制台",
                  "index": "doc2-8",
                  "path": "",
                  "pathRaw": "",
                  "children": [
                    {
                      "name": "控制台",
                      "index": "doc2-8-1",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F8%20%E6%8E%A7%E5%88%B6%E5%8F%B0%2F%E6%8E%A7%E5%88%B6%E5%8F%B0.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/8 控制台/控制台.md",
                      "children": null
                    }
                  ]
                },
                {
                  "name": "9 统计看板",
                  "index": "doc2-9",
                  "path": "",
                  "pathRaw": "",
                  "children": [
                    {
                      "name": "API调用统计",
                      "index": "doc2-9-1",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F9%20%E7%BB%9F%E8%AE%A1%E7%9C%8B%E6%9D%BF%2FAPI%E8%B0%83%E7%94%A8%E7%BB%9F%E8%AE%A1.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/9 统计看板/API调用统计.md",
                      "children": null
                    },
                    {
                      "name": "API调用记录",
                      "index": "doc2-9-2",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F9%20%E7%BB%9F%E8%AE%A1%E7%9C%8B%E6%9D%BF%2FAPI%E8%B0%83%E7%94%A8%E8%AE%B0%E5%BD%95.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/9 统计看板/API调用记录.md",
                      "children": null
                    }
                  ]
                },
                {
                  "name": "10 安全工具链",
                  "index": "doc2-10",
                  "path": "",
                  "pathRaw": "",
                  "children": [
                    {
                      "name": "敏感词管理",
                      "index": "doc2-10-1",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F10%20%E5%AE%89%E5%85%A8%E5%B7%A5%E5%85%B7%E9%93%BE%2F%E6%95%8F%E6%84%9F%E8%AF%8D%E7%AE%A1%E7%90%86.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/10 安全工具链/敏感词管理.md",
                      "children": null
                    }
                  ]
                },
                {
                  "name": "11 通用组件",
                  "index": "doc2-11",
                  "path": "",
                  "pathRaw": "",
                  "children": [
                    {
                      "name": "ChatConsult 专业咨询",
                      "index": "doc2-11-1",
                      "path": "v2.4.3-1747119948649%2F2%20%E5%B9%B3%E5%8F%B0%E6%93%8D%E4%BD%9C%E8%AF%B4%E6%98%8E%2F11%20%E9%80%9A%E7%94%A8%E7%BB%84%E4%BB%B6%2FChatConsult%20%E4%B8%93%E4%B8%9A%E5%92%A8%E8%AF%A2.md",
                      "pathRaw": "v2.4.3-1747119948649/2 平台操作说明/11 通用组件/ChatConsult 专业咨询.md",
                      "children": null
                    }
                  ]
                }
              ]
            },
            {
              "name": "3 开发者指南",
              "index": "doc3",
              "path": "",
              "pathRaw": "",
              "children": [
                {
                  "name": "API调用指南",
                  "index": "doc3-1",
                  "path": "v2.4.3-1747119948649%2F3%20%E5%BC%80%E5%8F%91%E8%80%85%E6%8C%87%E5%8D%97%2FAPI%E8%B0%83%E7%94%A8%E6%8C%87%E5%8D%97.md",
                  "pathRaw": "v2.4.3-1747119948649/3 开发者指南/API调用指南.md",
                  "children": null
                },
                {
                  "name": "获取access_token",
                  "index": "doc3-2",
                  "path": "v2.4.3-1747119948649%2F3%20%E5%BC%80%E5%8F%91%E8%80%85%E6%8C%87%E5%8D%97%2F%E8%8E%B7%E5%8F%96access_token.md",
                  "pathRaw": "v2.4.3-1747119948649/3 开发者指南/获取access_token.md",
                  "children": null
                }
              ]
            },
            {
              "name": "4 大语言模型",
              "index": "doc4",
              "path": "",
              "pathRaw": "",
              "children": [
                {
                  "name": "大语言模型",
                  "index": "doc4-1",
                  "path": "",
                  "pathRaw": "",
                  "children": [
                    {
                      "name": "DeepSeek-R1-671B",
                      "index": "doc4-1-1",
                      "path": "v2.4.3-1747119948649%2F4%20%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2F%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2FDeepSeek-R1-671B.md",
                      "pathRaw": "v2.4.3-1747119948649/4 大语言模型/大语言模型/DeepSeek-R1-671B.md",
                      "children": null
                    },
                    {
                      "name": "DeepSeek-R1-Distill-Qwen-1.5B",
                      "index": "doc4-1-2",
                      "path": "v2.4.3-1747119948649%2F4%20%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2F%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2FDeepSeek-R1-Distill-Qwen-1.5B.md",
                      "pathRaw": "v2.4.3-1747119948649/4 大语言模型/大语言模型/DeepSeek-R1-Distill-Qwen-1.5B.md",
                      "children": null
                    },
                    {
                      "name": "DeepSeek-R1-Distill-Qwen-14B",
                      "index": "doc4-1-3",
                      "path": "v2.4.3-1747119948649%2F4%20%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2F%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2FDeepSeek-R1-Distill-Qwen-14B.md",
                      "pathRaw": "v2.4.3-1747119948649/4 大语言模型/大语言模型/DeepSeek-R1-Distill-Qwen-14B.md",
                      "children": null
                    },
                    {
                      "name": "DeepSeek-R1-Distill-Qwen-32B",
                      "index": "doc4-1-4",
                      "path": "v2.4.3-1747119948649%2F4%20%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2F%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2FDeepSeek-R1-Distill-Qwen-32B.md",
                      "pathRaw": "v2.4.3-1747119948649/4 大语言模型/大语言模型/DeepSeek-R1-Distill-Qwen-32B.md",
                      "children": null
                    },
                    {
                      "name": "DeepSeek-R1-Distill-Qwen-7B",
                      "index": "doc4-1-5",
                      "path": "v2.4.3-1747119948649%2F4%20%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2F%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2FDeepSeek-R1-Distill-Qwen-7B.md",
                      "pathRaw": "v2.4.3-1747119948649/4 大语言模型/大语言模型/DeepSeek-R1-Distill-Qwen-7B.md",
                      "children": null
                    },
                    {
                      "name": "DeepSeek-V3-671B",
                      "index": "doc4-1-6",
                      "path": "v2.4.3-1747119948649%2F4%20%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2F%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2FDeepSeek-V3-671B.md",
                      "pathRaw": "v2.4.3-1747119948649/4 大语言模型/大语言模型/DeepSeek-V3-671B.md",
                      "children": null
                    },
                    {
                      "name": "Unicom-13B-Chat",
                      "index": "doc4-1-7",
                      "path": "v2.4.3-1747119948649%2F4%20%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2F%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2FUnicom-13B-Chat.md",
                      "pathRaw": "v2.4.3-1747119948649/4 大语言模型/大语言模型/Unicom-13B-Chat.md",
                      "children": null
                    },
                    {
                      "name": "Unicom-34B-Chat",
                      "index": "doc4-1-8",
                      "path": "v2.4.3-1747119948649%2F4%20%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2F%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2FUnicom-34B-Chat.md",
                      "pathRaw": "v2.4.3-1747119948649/4 大语言模型/大语言模型/Unicom-34B-Chat.md",
                      "children": null
                    },
                    {
                      "name": "Unicom-34B-Coder",
                      "index": "doc4-1-9",
                      "path": "v2.4.3-1747119948649%2F4%20%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2F%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2FUnicom-34B-Coder.md",
                      "pathRaw": "v2.4.3-1747119948649/4 大语言模型/大语言模型/Unicom-34B-Coder.md",
                      "children": null
                    },
                    {
                      "name": "Unicom-70B-Chat",
                      "index": "doc4-1-10",
                      "path": "v2.4.3-1747119948649%2F4%20%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2F%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2FUnicom-70B-Chat.md",
                      "pathRaw": "v2.4.3-1747119948649/4 大语言模型/大语言模型/Unicom-70B-Chat.md",
                      "children": null
                    },
                    {
                      "name": "Unicom-7B-Chat",
                      "index": "doc4-1-11",
                      "path": "v2.4.3-1747119948649%2F4%20%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2F%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2FUnicom-7B-Chat.md",
                      "pathRaw": "v2.4.3-1747119948649/4 大语言模型/大语言模型/Unicom-7B-Chat.md",
                      "children": null
                    }
                  ]
                },
                {
                  "name": "推理服务API列表",
                  "index": "doc4-2",
                  "path": "v2.4.3-1747119948649%2F4%20%E5%A4%A7%E8%AF%AD%E8%A8%80%E6%A8%A1%E5%9E%8B%2F%E6%8E%A8%E7%90%86%E6%9C%8D%E5%8A%A1API%E5%88%97%E8%A1%A8.md",
                  "pathRaw": "v2.4.3-1747119948649/4 大语言模型/推理服务API列表.md",
                  "children": null
                }
              ]
            },
            {
              "name": "5 知识库管理API",
              "index": "doc5",
              "path": "",
              "pathRaw": "",
              "children": [
                {
                  "name": "上传文件并异步向量化",
                  "index": "doc5-1",
                  "path": "v2.4.3-1747119948649%2F5%20%E7%9F%A5%E8%AF%86%E5%BA%93%E7%AE%A1%E7%90%86API%2F%E4%B8%8A%E4%BC%A0%E6%96%87%E4%BB%B6%E5%B9%B6%E5%BC%82%E6%AD%A5%E5%90%91%E9%87%8F%E5%8C%96.md",
                  "pathRaw": "v2.4.3-1747119948649/5 知识库管理API/上传文件并异步向量化.md",
                  "children": null
                },
                {
                  "name": "创建知识库",
                  "index": "doc5-2",
                  "path": "v2.4.3-1747119948649%2F5%20%E7%9F%A5%E8%AF%86%E5%BA%93%E7%AE%A1%E7%90%86API%2F%E5%88%9B%E5%BB%BA%E7%9F%A5%E8%AF%86%E5%BA%93.md",
                  "pathRaw": "v2.4.3-1747119948649/5 知识库管理API/创建知识库.md",
                  "children": null
                },
                {
                  "name": "创建知识库（多参数）",
                  "index": "doc5-3",
                  "path": "v2.4.3-1747119948649%2F5%20%E7%9F%A5%E8%AF%86%E5%BA%93%E7%AE%A1%E7%90%86API%2F%E5%88%9B%E5%BB%BA%E7%9F%A5%E8%AF%86%E5%BA%93%EF%BC%88%E5%A4%9A%E5%8F%82%E6%95%B0%EF%BC%89.md",
                  "pathRaw": "v2.4.3-1747119948649/5 知识库管理API/创建知识库（多参数）.md",
                  "children": null
                },
                {
                  "name": "删除知识库",
                  "index": "doc5-4",
                  "path": "v2.4.3-1747119948649%2F5%20%E7%9F%A5%E8%AF%86%E5%BA%93%E7%AE%A1%E7%90%86API%2F%E5%88%A0%E9%99%A4%E7%9F%A5%E8%AF%86%E5%BA%93.md",
                  "pathRaw": "v2.4.3-1747119948649/5 知识库管理API/删除知识库.md",
                  "children": null
                },
                {
                  "name": "删除知识库内文件",
                  "index": "doc5-5",
                  "path": "v2.4.3-1747119948649%2F5%20%E7%9F%A5%E8%AF%86%E5%BA%93%E7%AE%A1%E7%90%86API%2F%E5%88%A0%E9%99%A4%E7%9F%A5%E8%AF%86%E5%BA%93%E5%86%85%E6%96%87%E4%BB%B6.md",
                  "pathRaw": "v2.4.3-1747119948649/5 知识库管理API/删除知识库内文件.md",
                  "children": null
                },
                {
                  "name": "知识库列表",
                  "index": "doc5-6",
                  "path": "v2.4.3-1747119948649%2F5%20%E7%9F%A5%E8%AF%86%E5%BA%93%E7%AE%A1%E7%90%86API%2F%E7%9F%A5%E8%AF%86%E5%BA%93%E5%88%97%E8%A1%A8.md",
                  "pathRaw": "v2.4.3-1747119948649/5 知识库管理API/知识库列表.md",
                  "children": null
                },
                {
                  "name": "知识库列表（多参数）",
                  "index": "doc5-7",
                  "path": "v2.4.3-1747119948649%2F5%20%E7%9F%A5%E8%AF%86%E5%BA%93%E7%AE%A1%E7%90%86API%2F%E7%9F%A5%E8%AF%86%E5%BA%93%E5%88%97%E8%A1%A8%EF%BC%88%E5%A4%9A%E5%8F%82%E6%95%B0%EF%BC%89.md",
                  "pathRaw": "v2.4.3-1747119948649/5 知识库管理API/知识库列表（多参数）.md",
                  "children": null
                },
                {
                  "name": "获取文件列表",
                  "index": "doc5-8",
                  "path": "v2.4.3-1747119948649%2F5%20%E7%9F%A5%E8%AF%86%E5%BA%93%E7%AE%A1%E7%90%86API%2F%E8%8E%B7%E5%8F%96%E6%96%87%E4%BB%B6%E5%88%97%E8%A1%A8.md",
                  "pathRaw": "v2.4.3-1747119948649/5 知识库管理API/获取文件列表.md",
                  "children": null
                },
                {
                  "name": "获取最近上传的知识",
                  "index": "doc5-9",
                  "path": "v2.4.3-1747119948649%2F5%20%E7%9F%A5%E8%AF%86%E5%BA%93%E7%AE%A1%E7%90%86API%2F%E8%8E%B7%E5%8F%96%E6%9C%80%E8%BF%91%E4%B8%8A%E4%BC%A0%E7%9A%84%E7%9F%A5%E8%AF%86.md",
                  "pathRaw": "v2.4.3-1747119948649/5 知识库管理API/获取最近上传的知识.md",
                  "children": null
                }
              ]
            },
            {
              "name": "6 智能体能力拓展",
              "index": "doc6",
              "path": "",
              "pathRaw": "",
              "children": [
                {
                  "name": "Minio文件上传API",
                  "index": "doc6-1",
                  "path": "",
                  "pathRaw": "",
                  "children": [
                    {
                      "name": "minio文件上传",
                      "index": "doc6-1-1",
                      "path": "v2.4.3-1747119948649%2F6%20%E6%99%BA%E8%83%BD%E4%BD%93%E8%83%BD%E5%8A%9B%E6%8B%93%E5%B1%95%2FMinio%E6%96%87%E4%BB%B6%E4%B8%8A%E4%BC%A0API%2Fminio%E6%96%87%E4%BB%B6%E4%B8%8A%E4%BC%A0.md",
                      "pathRaw": "v2.4.3-1747119948649/6 智能体能力拓展/Minio文件上传API/minio文件上传.md",
                      "children": null
                    }
                  ]
                },
                {
                  "name": "RAG API",
                  "index": "doc6-2",
                  "path": "",
                  "pathRaw": "",
                  "children": [
                    {
                      "name": "单机版-知识增强流式返回",
                      "index": "doc6-2-1",
                      "path": "v2.4.3-1747119948649%2F6%20%E6%99%BA%E8%83%BD%E4%BD%93%E8%83%BD%E5%8A%9B%E6%8B%93%E5%B1%95%2FRAG%20API%2F%E5%8D%95%E6%9C%BA%E7%89%88-%E7%9F%A5%E8%AF%86%E5%A2%9E%E5%BC%BA%E6%B5%81%E5%BC%8F%E8%BF%94%E5%9B%9E.md",
                      "pathRaw": "v2.4.3-1747119948649/6 智能体能力拓展/RAG API/单机版-知识增强流式返回.md",
                      "children": null
                    }
                  ]
                },
                {
                  "name": "智能体 API",
                  "index": "doc6-3",
                  "path": "",
                  "pathRaw": "",
                  "children": [
                    {
                      "name": "智能体",
                      "index": "doc6-3-1",
                      "path": "v2.4.3-1747119948649%2F6%20%E6%99%BA%E8%83%BD%E4%BD%93%E8%83%BD%E5%8A%9B%E6%8B%93%E5%B1%95%2F%E6%99%BA%E8%83%BD%E4%BD%93%20API%2F%E6%99%BA%E8%83%BD%E4%BD%93.md",
                      "pathRaw": "v2.4.3-1747119948649/6 智能体能力拓展/智能体 API/智能体.md",
                      "children": null
                    }
                  ]
                }
              ]
            },
            {
              "name": "7 安全工具链API",
              "index": "doc7",
              "path": "",
              "pathRaw": "",
              "children": [
                {
                  "name": "敏感词检查接口服务",
                  "index": "doc7-1",
                  "path": "v2.4.3-1747119948649%2F7%20%E5%AE%89%E5%85%A8%E5%B7%A5%E5%85%B7%E9%93%BEAPI%2F%E6%95%8F%E6%84%9F%E8%AF%8D%E6%A3%80%E6%9F%A5%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1.md",
                  "pathRaw": "v2.4.3-1747119948649/7 安全工具链API/敏感词检查接口服务.md",
                  "children": null
                }
              ]
            }
          ],
          "msg": ""
        }
        this.docMenuList = data || []
      })
    },
    getCurrentMenu() {
      const route = this.$route
      this.getDocMenu().then(() => {
        this.getMenuList(route)
      })
    },
    getMenuList() {
      const { params, path } = this.$route || {}
      const { id } = params || {}
      let val = path
      // 获取当前菜单列表
      const menus = this.docMenuList
      this.menuList = menus
      this.defaultOpeneds = menus.map(item => item.index)

      // 跳转到文档中心第一个菜单栏
      if (id === DOC_FIRST_KEY) {
        const { path } = fetchPermFirPath(menus)
        val = path
        this.$router.push({path})
      }

      // 给当前 activeIndex 赋值
      this.changeMenuIndex(fetchCurrentPathIndex(val, menus))
    },
    changeMenuIndex(index) {
      this.activeIndex = index
    },
  }
}
</script>

<style lang="scss" scoped>
.doc-page-container {
  height:100%;
  .doc-outer-container{
    height: calc(100% - 90px);
    background: #fff;
    width: calc(100% - 130px);
    margin: 0 auto;
    border-radius: 8px;
    /*element ui 样式重写*/
    .doc-inner-container{
      height: 100%;
      display: flex;
      .el-aside.full-menu-aside {
        height: 100%;
        background-color: #fff;
        overflow-y: auto;
        overflow-x: auto;
        border-right: 1px solid #ededed;
        border-radius: 8px 0 0 0;
        .el-menu{
          min-height: 100%;
          width: fit-content;
          overflow-x: auto;
          overflow-y: hidden;
          .menu-indent /deep/ .el-submenu__title,
          .menu-indent-item {
            padding-left: 49px !important;
          }
          .menu-indent-sub /deep/ .el-submenu__title{
            padding-left: 60px !important;
          }
        }
      }
      .doc-page-main {
        width: 100%;
        height: 100%;
        min-height: 580px;
        overflow: auto;
        padding-top: 30px;
        background: rgba(255, 255, 255, 0);
        border-radius: 8px;
      }
      /deep/ .el-submenu__title,
      /deep/ .el-menu-item span,
      /deep/ .el-submenu__title span {
        font-size: 14px !important;
      }
      /deep/ .el-menu-item.is-active {
        background-color: rgba(230, 0, 1, 0.03) !important;
      }
      /deep/ .el-menu-item:focus {
        background-color: rgba(230, 0, 1, 0.03) !important;
      }
      /deep/ .el-menu-item:hover,
      /deep/ .el-submenu__title:hover {
        background-color: #f6f6f6 !important;
      }
      /deep/ .el-submenu__title {
        span {
          font-size: 14px !important;
        }
      }
      /deep/ .el-submenu.is-active .el-submenu__title {
        border-bottom-color: $color !important;
      }
    }
  }
  .doc-outer-container /deep/ {
    .el-submenu.is-active,
    .el-submenu.is-active > .el-submenu__title,
    .el-submenu.is-active > .el-submenu__title i:first-child,
    .el-submenu.is-active > .el-submenu__title .el-submenu__icon-arrow {
      color: $color !important;
    }
  }
}
.doc-header{
  padding: 20px 0;
  display: flex;
  justify-content: center;
  .doc-header__left {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 30%;
    span {
      font-size: 22px;
      margin-left: 18px;
      font-weight: bold;
      color: $color_title;
    }
  }
  .doc-header__right {
    width: 60%;
    display: flex;
    align-items: center;
  }
  .top-search-input /deep/ {
    .el-input__inner {
      border-radius: 20px;
    }
    .top-search-popover {
      left: auto !important;
      right: 10px !important;
    }

    .popper__arrow {
      display: none !important;
    }
    .el-select-dropdown__wrap {
      max-height: 550px !important;
    }
    .el-select-dropdown__item {
      background: rgba(255, 255, 255, 0) !important;
      padding: 4px 10px;
    }
    .el-input__suffix-inner {
      display: inline-block;
    }
    .header_search-option {
      width: 500px;
      .header_search-title {
        color: #fff;
        background-color: #df0000;
        padding: 0 10px;
        font-weight: bold;
        font-size: 16px;
      }
      .header_search-item {
        display: flex;
        justify-content: space-between;
        color: #333;
        border-bottom: 1px solid #d8d6d6;
        .header_search-item-left {
          background-color: #f1f1f1;
          width: 40%;
          text-align: right;
          padding: 0 10px;
          font-weight: bold;
        }
        .header_search-item-right {
          width: 60%;
          text-align: left;
          padding: 3px 10px;
          color: #666;
          font-weight: bold;
          white-space: normal; /* 保留空白符序列，但是正常换行 */
          word-break: break-all;
          span {
            color: #df0000;
            text-decoration: underline;
          }
          p {
            line-height: 20px;
          }
          a {
            color: #df0000;
          }
          img {
            width: 100%;
          }
        }
        .header_search-item-right:hover {
          background-color: #f1f1f1;
        }
      }
    }
  }
}
</style>
