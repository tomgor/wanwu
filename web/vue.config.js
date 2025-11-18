"use strict";

const path = require("path");
const { name } = require('./package.json')
const webpack = require('webpack')
const VersionInfoPlugin = require('./build/version-plugin');
function resolve(dir) {
  return path.join(__dirname, dir);
}

const CompressionWebpackPlugin = require('compression-webpack-plugin')
const isProdOrTest = process.env.NODE_ENV !== 'development'

module.exports = {
  // 基础配置 详情看文档
  publicPath: '/',
  outputDir: "dist",
  assetsDir: "static",
  lintOnSave: process.env.NODE_ENV === "development",
  productionSourceMap: false,//源码映射
  chainWebpack(config) {
    config.module
      .rule('md')
      .test(/\.md$/)
      .use('html-loader')
      .loader('html-loader')
      .end()
      .use('markdown-loader')
      .loader('markdown-loader')
      .end()

    config.plugins.delete('prefetch')
    if (isProdOrTest) {
      // 对超过10kb的文件gzip压缩
      config.plugin('compressionPlugin').use(
        new CompressionWebpackPlugin({
          test: /\.(css|html)$/,
          threshold: 10240,
        })
      )
    }
    config.module
      .rule('svg')
      .exclude.add(resolve('src/assets/icons'))	//svg文件位置
      .end()
    config.module
      .rule('icons')
      .test(/\.svg$/)
      .include.add(resolve('src/assets/icons'))	//svg文件位置
      .end()
      .use('svg-sprite-loader')
      .loader('svg-sprite-loader')
      .options({
        symbolId: 'icon-[name]'
      })
      .end()
  },
  devServer: {
    port: 8080,
    open: false,
    hot: true,
    overlay: {
      warnings: false,
      errors: true,
    },
    publicPath: '/',
    headers: {
      'Access-Control-Allow-Origin': '*',
    },
    proxy: {
      '/openAi': {
        target: process.env.VUE_BASE_URL, changeOrigin: true, secure: false, bypass(req, res, options) {
          const realUrl = options.target + (options.rewrite ? options.rewrite(req.url) : '')
          console.log(realUrl) // 在终端显示
          res.setHeader('A-Real-Url', realUrl) // 添加响应标头(A-Real-Url为自定义命名)，在浏览器中显示
        },
      },
      '/api': {
        target: process.env.VUE_BASE_URL, changeOrigin: true, secure: false, bypass(req, res, options) {
          const realUrl = options.target + (options.rewrite ? options.rewrite(req.url) : '')
          console.log(realUrl) // 在终端显示
          res.setHeader('A-Real-Url', realUrl) // 添加响应标头(A-Real-Url为自定义命名)，在浏览器中显示
        },
      },
      '/workflow/api': {
        target: process.env.VUE_BASE_URL, changeOrigin: true, secure: false, bypass(req, res, options) {
          const realUrl = options.target + (options.rewrite ? options.rewrite(req.url) : '')
          console.log(realUrl) // 在终端显示
          res.setHeader('A-Real-Url', realUrl) // 添加响应标头(A-Real-Url为自定义命名)，在浏览器中显示
        },
      },
      '/user/api': {
        target: process.env.VUE_BASE_URL, changeOrigin: true, secure: false, bypass(req, res, options) {
          const realUrl = options.target + (options.rewrite ? options.rewrite(req.url) : '')
          console.log(realUrl) // 在终端显示
          res.setHeader('A-Real-Url', realUrl) // 添加响应标头(A-Real-Url为自定义命名)，在浏览器中显示
        },
      },
      '/service/url/openurl/v1': {
        target: process.env.VUE_BASE_URL, changeOrigin: true, secure: false, bypass(req, res, options) {
          const realUrl = options.target + (options.rewrite ? options.rewrite(req.url) : '')
          console.log(realUrl) // 在终端显示
          res.setHeader('A-Real-Url', realUrl) // 添加响应标头(A-Real-Url为自定义命名)，在浏览器中显示
        },
      },
      '/service/api': {
        target: process.env.VUE_BASE_URL, changeOrigin: true, secure: false, bypass(req, res, options) {
          const realUrl = options.target + (options.rewrite ? options.rewrite(req.url) : '')
          console.log(realUrl) // 在终端显示
          res.setHeader('A-Real-Url', realUrl) // 添加响应标头(A-Real-Url为自定义命名)，在浏览器中显示
        },
      },
      '/training/api': {
        target: process.env.VUE_BASE_URL, changeOrigin: true, secure: false, bypass(req, res, options) {
          const realUrl = options.target + (options.rewrite ? options.rewrite(req.url) : '')
          console.log(realUrl) // 在终端显示
          res.setHeader('A-Real-Url', realUrl) // 添加响应标头(A-Real-Url为自定义命名)，在浏览器中显示
        },
      },
      '/resource/api': {
        target: process.env.VUE_BASE_URL, changeOrigin: true, secure: false, bypass(req, res, options) {
          const realUrl = options.target + (options.rewrite ? options.rewrite(req.url) : '')
          console.log(realUrl) // 在终端显示
          res.setHeader('A-Real-Url', realUrl) // 添加响应标头(A-Real-Url为自定义命名)，在浏览器中显示
        },
      },
      '/datacenter/api': {
        target: process.env.VUE_BASE_URL, changeOrigin: true, secure: false, bypass(req, res, options) {
          const realUrl = options.target + (options.rewrite ? options.rewrite(req.url) : '')
          console.log(realUrl) // 在终端显示
          res.setHeader('A-Real-Url', realUrl) // 添加响应标头(A-Real-Url为自定义命名)，在浏览器中显示
        },
      },
      '/modelprocess/api': {
        target: process.env.VUE_BASE_URL, changeOrigin: true, secure: false, bypass(req, res, options) {
          const realUrl = options.target + (options.rewrite ? options.rewrite(req.url) : '')
          console.log(realUrl) // 在终端显示
          res.setHeader('A-Real-Url', realUrl) // 添加响应标头(A-Real-Url为自定义命名)，在浏览器中显示
        },
      },
      '/expand/api': {
        target: process.env.VUE_BASE_URL, changeOrigin: true, secure: false, bypass(req, res, options) {
          const realUrl = options.target + (options.rewrite ? options.rewrite(req.url) : '')
          console.log(realUrl) // 在终端显示
          res.setHeader('A-Real-Url', realUrl) // 添加响应标头(A-Real-Url为自定义命名)，在浏览器中显示
        },
      },
      '/record/api': {
        target: process.env.VUE_BASE_URL, changeOrigin: true, secure: false, bypass(req, res, options) {
          const realUrl = options.target + (options.rewrite ? options.rewrite(req.url) : '')
          console.log(realUrl) // 在终端显示
          res.setHeader('A-Real-Url', realUrl) // 添加响应标头(A-Real-Url为自定义命名)，在浏览器中显示
        },
      },
      '/img': {
        target: process.env.VUE_BASE_URL, changeOrigin: true, secure: false, bypass(req, res, options) {
          const realUrl = options.target + (options.rewrite ? options.rewrite(req.url) : '')
          console.log(realUrl) // 在终端显示
          res.setHeader('A-Real-Url', realUrl) // 添加响应标头(A-Real-Url为自定义命名)，在浏览器中显示
        },
      },
      '/konwledgeServe': {
        target: process.env.VUE_BASE_URL, changeOrigin: true, secure: false, bypass(req, res, options) {
          const realUrl = options.target + (options.rewrite ? options.rewrite(req.url) : '')
          console.log(realUrl) // 在终端显示
          res.setHeader('A-Real-Url', realUrl) // 添加响应标头(A-Real-Url为自定义命名)，在浏览器中显示
        },
      },
      '/proxyupload': {
        target: process.env.VUE_BASE_URL, changeOrigin: true, secure: false, bypass(req, res, options) {
          const realUrl = options.target + (options.rewrite ? options.rewrite(req.url) : '')
          console.log(realUrl) // 在终端显示
          res.setHeader('A-Real-Url', realUrl) // 添加响应标头(A-Real-Url为自定义命名)，在浏览器中显示
        },
      },
      '/use/model/api': {
        target: process.env.VUE_BASE_URL, changeOrigin: true, secure: false, bypass(req, res, options) {
          const realUrl = options.target + (options.rewrite ? options.rewrite(req.url) : '')
          console.log(realUrl) // 在终端显示
          res.setHeader('A-Real-Url', realUrl) // 添加响应标头(A-Real-Url为自定义命名)，在浏览器中显示
        },
      },
      '/prompt/api': {
        target: process.env.VUE_BASE_URL, changeOrigin: true, secure: false, bypass(req, res, options) {
          const realUrl = options.target + (options.rewrite ? options.rewrite(req.url) : '')
          console.log(realUrl) // 在终端显示
          res.setHeader('A-Real-Url', realUrl) // 添加响应标头(A-Real-Url为自定义命名)，在浏览器中显示
        },
      },
    },
  },
  css: {
    sourceMap: false,
    loaderOptions: {
      sass: {
        prependData: `@import "~@/style/theme/vars_blue.scss";@import "~@/style/theme/common.scss";` // 假设variables.scss位于src/styles目录下
      }
    }
  },
  configureWebpack: {
    //    @路径走src文件夹
    module: {
      rules: [
        {
          test: /\.swf$/,
          loader: "url-loader",
          options: {
            limit: 10000,
            name: "static/media/[name].[hash:7].[ext]",
          },
        },
      ],
    },
    resolve: {
      alias: {
        'vue$': 'vue/dist/vue.esm.js',
        "@": resolve("src"),
        "@common": resolve("common"),
      },
    },
    output: {
      // 把子应用打包成 umd 库格式(必须)
      library: `${name}-[name]`,
      libraryTarget: 'umd',
      jsonpFunction: `webpackJsonp_${name}`,
    },
    plugins: [
      new webpack.optimize.LimitChunkCountPlugin({
        maxChunks: 10, // 来限制 chunk 的最大数量
      }),
      new webpack.optimize.MinChunkSizePlugin({
        minChunkSize: 50000 // Minimum number of characters
      }),
      new VersionInfoPlugin()
    ]
  },
};
