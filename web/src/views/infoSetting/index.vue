<template>
  <div>
    <!--tab配置-->
    <el-card shadow="never" class="docPage-card">
      <div slot="header">
        <span class="card-title">{{$t('infoSetting.tabSet')}}</span>
      </div>
      <el-form label-width="120px" :model="tabForm" :rules="titleRules" ref="tabForm">
        <el-form-item :label="$t('infoSetting.form.labelTitle')" prop="tabTitle">
          <el-input v-model="tabForm.tabTitle" style="width: 300px" />
          <div style="font-size: 11px; color: #aaa; margin-top: -9px">
            {{$t('infoSetting.hint.labelTitle')}}
          </div>
        </el-form-item>
        <el-form-item :label="$t('infoSetting.form.labelIcon')" prop="tabLogo">
          <el-upload
            class="avatar-uploader"
            action=""
            name="files"
            :show-file-list="false"
            :http-request="handleUploadLabelIcon"
            :on-error="handleUploadError"
            accept=".png,.jpg,.jpeg"
          >
            <img v-if="tabLogo" :src="tabLogo" class="avatar">
            <i v-if="!tabLogo" class="el-icon-plus avatar-uploader-icon"></i>
            <span style="margin-left: 12px; color: #aaa !important;">
              {{$t('infoSetting.hint.imgUpload')}}
            </span>
            <div style="text-align: left; font-size: 11px; color: #aaa; margin-top: -9px">
              {{$t('infoSetting.hint.labelIcon')}}
            </div>
          </el-upload>
        </el-form-item>
      </el-form>
      <div class="card-footer">
        <el-button
          :loading="titleLoading"
          type="primary"
          @click="handleSubmitTab()"
        >
          {{$t('infoSetting.form.save')}}
        </el-button>
      </div>
    </el-card>
    <!--登录页配置-->
    <el-card shadow="never" class="docPage-card">
      <div slot="header">
        <span class="card-title">{{$t('infoSetting.loginBgSet')}}</span>
      </div>
      <el-form label-width="120px" :model="loginForm" :rules="loginRules" ref="loginForm">
        <el-form-item :label="$t('infoSetting.form.loginBg')" prop="loginBg">
          <el-upload
            class="avatar-uploader"
            action=""
            name="files"
            :show-file-list="false"
            :http-request="handleUploadLoginBg"
            :on-error="handleUploadError"
            accept=".png,.jpg,.jpeg"
          >
            <img v-if="loginBg" :src="loginBg" class="avatar">
            <i v-if="!loginBg" class="el-icon-plus avatar-uploader-icon"></i>
            <span style="margin-left: 12px; color: #aaa !important;">
              {{$t('infoSetting.hint.imgUpload')}}
            </span>
            <div style="text-align: left; font-size: 11px; color: #aaa; margin-top: -9px">
              {{$t('infoSetting.hint.loginBg')}}
            </div>
          </el-upload>
        </el-form-item>
        <el-form-item :label="$t('infoSetting.form.logoWelcome')" prop="logoWelcomeText">
          <el-input v-model="loginForm.loginWelcomeText" style="width: 300px" />
          <div style="font-size: 11px; color: #aaa; margin-top: -9px">
            {{$t('infoSetting.hint.logoWelcome')}}
          </div>
        </el-form-item>
        <el-form-item :label="$t('infoSetting.form.loginButtonColor')" prop="loginButtonColor">
          <el-color-picker v-model="loginForm.loginButtonColor" show-alpha></el-color-picker>
          <div style="text-align: left; font-size: 11px; color: #aaa; margin-top: -9px">
            {{$t('infoSetting.hint.loginButtonColor')}}
          </div>
        </el-form-item>
      </el-form>
      <div class="card-footer">
        <el-button
          :loading="backgroundLoading"
          type="primary"
          @click="handelSubmitLogin()"
        >
          {{$t('infoSetting.form.save')}}
        </el-button>
      </div>
    </el-card>
    <!--平台配置-->
    <el-card shadow="never" class="docPage-card">
      <div slot="header">
        <span class="card-title">{{$t('infoSetting.platformSet')}}</span>
      </div>
      <el-form label-width="120px" :model="form" :rules="rules" ref="form">
        <el-form-item :label="$t('infoSetting.form.platformTitle')" prop="platformTitle">
          <el-input v-model="form.platformTitle" style="width: 300px" />
          <div style="font-size: 11px; color: #aaa; margin-top: -9px">
            {{$t('infoSetting.hint.platformTitle')}}
          </div>
        </el-form-item>
        <el-form-item :label="$t('infoSetting.form.logo')" prop="platformLogo">
          <el-upload
            class="avatar-uploader"
            action=""
            name="files"
            :show-file-list="false"
            :http-request="handleUploadLogo"
            :on-error="handleUploadError"
            accept=".png,.jpg,.jpeg"
          >
            <img v-if="platformLogo" :src="platformLogo" class="avatar">
            <i v-if="!platformLogo" class="el-icon-plus avatar-uploader-icon"></i>
            <span style="margin-left: 12px; color: #aaa !important;">
              {{$t('infoSetting.hint.imgUpload')}}
            </span>
            <div style="text-align: left; font-size: 11px; color: #aaa; margin-top: -9px">
              {{$t('infoSetting.hint.logo')}}
            </div>
          </el-upload>
        </el-form-item>
        <el-form-item :label="$t('infoSetting.form.bgColor')" prop="color">
          <el-radio-group v-model="radio">
            <el-radio label="0">{{$t('infoSetting.hint.oneColor')}}</el-radio>
            <el-radio label="1">{{$t('infoSetting.hint.linearColor')}}</el-radio>
          </el-radio-group>
          <div v-if="radio === '0'">
            <el-color-picker v-model="form.color" show-alpha></el-color-picker>
          </div>
          <div v-if="radio === '1'">
            <el-input v-model="form.color" style="width: 300px" />
            <span :style="`display: inline-block; vertical-align: middle; width: 32px; height: 32px; border-radius: 2px; margin-left: 14px; background: ${form.color}`"></span>
          </div>
          <div style="text-align: left; font-size: 11px; color: #aaa; margin-top: -9px">
            {{$t('infoSetting.hint.bgColor')}}
          </div>
        </el-form-item>
      </el-form>
      <div class="card-footer">
        <el-button
          :loading="logoLoading"
          type="primary"
          @click="handleSubmit()"
        >
          {{$t('infoSetting.form.save')}}
        </el-button>
      </div>
    </el-card>
  </div>
</template>
<script>

import { setPlatformInfo } from "@/api/setInfo"
import { mapActions, mapGetters } from "vuex"

export default {
  data() {
    const checkImage = (val, rule, value, callback) => {
      // 判断编辑器值，编辑器默认没输入时是 <p><br></p>
      if (!val) {
        callback(new Error(this.$t('infoSetting.form.uploadHint')))
      } else {
        return callback()
      }
    }
    return {
      titleLoading: false,
      logoLoading: false,
      backgroundLoading: false,
      tabForm: {
        tabTitle: '',
        tabLogo: '',
      },
      radio: '0',
      form: {
        platformTitle: '',
        platformLogo: '',
        color: ''
      },
      loginForm: {
        loginBg: '',
        loginButtonColor: '',
        loginWelcomeText: ''
      },
      tabLogo: '',
      platformLogo: '',
      loginBg: '',
      tabRules: {
        tabTitle: [{required: true, message: this.$t('common.input.placeholder'), trigger: 'blur'}],
        tabLogo: [{required: true, validator: (...params) => checkImage(this.tabLogo, ...params), trigger: 'change'}],
      },
      rules: {
        platformTitle: [{required: true, message: this.$t('common.input.placeholder'), trigger: 'blur'}],
        platformLogo: [{required: true, validator: (...params) => checkImage(this.platformLogo, ...params), trigger: 'change'}],
        color: [{required: true, message: this.$t('common.select.placeholder'), trigger: 'change'}],
      },
      loginRules: {
        loginBg: [{required: true, validator: (...params) => checkImage(this.loginBg, ...params), trigger: 'change'}],
        loginButtonColor: [{required: true, message: this.$t('common.select.placeholder'), trigger: 'change'}],
        logoWelcomeText: [{required: true, message: this.$t('common.input.placeholder'), trigger: 'blur'}],
      }
    }
  },
  created() {
    this.getCommonInfo()
  },
  watch: {
    commonInfo:{
      handler(val) {
        const { home = {}, tab = {}, login = {} } = val ? val.data || {} : {}
        console.log(val, '----------------info')

        this.tabForm.tabTitle = tab.title
        this.tabForm.tabLogo = tab.logoPath
        this.tabLogo = tab.logoPath ? this.$basePath + '/user/api' + tab.logoPath : ''
        
        this.form.platformTitle = home.title
        this.form.platformLogo = home.logoPath
        this.platformLogo = home.logoPath ? this.$basePath + '/user/api' + home.logoPath : ''
        const color = home.navigationColor || 'linear-gradient(1deg, #FFFFFF 42%, #FFFFFF 42%, #EBEDFE 98%, #EEF0FF 98%)'
        this.radio = color.includes('linear-gradient') ? '1' : '0'
        this.form.color = color

        this.loginForm.loginBg = login.backgroundPath
        this.loginBg = login.backgroundPath ? this.$basePath + '/user/api' + login.backgroundPath : ''
        this.loginForm.loginWelcomeText = login.welcomeText || ''
        this.loginForm.loginButtonColor = login.loginButtonColor || '#d33a3a'
      },
      deep: true
    }
  },
  computed: {
    ...mapGetters('user', ['commonInfo']),
  },
  methods: {
    ...mapActions('user', ['getCommonInfo']),
    handleUploadLogo(data) {
      if (data.file) {
        this.form.platformLogo = data.file
        this.platformLogo = URL.createObjectURL(data.file)
      }
    },
    handleUploadLoginBg(data) {
      if (data.file) {
        this.loginForm.loginBg = data.file
        this.loginBg = URL.createObjectURL(data.file)
      }
    },
    handleUploadLabelIcon(data) {
      if (data.file) {
        this.tabForm.tabLogo = data.file
        this.tabLogo = URL.createObjectURL(data.file)
      }
    },
    handleUploadError() {
      this.$message.error(this.$t('common.message.uploadError'))
    },
    async submitData(type, data, isJson) {
      try {
        const config = isJson ? {} : {headers: { "Content-Type": "multipart/form-data" }}
        const formData = new FormData()
        for (let key in data) {
          formData.append(key, data[key])
        }
        this[`${type}Loading`] = true
        const res = await setPlatformInfo(type, isJson ? data : formData, config)
        if (res.code === 0) this.$message.success(this.$t('common.message.success'))
        if (type !== 'login') window.location.reload()
      } finally {
        this[`${type}Loading`] = false
      }
    },
    handleSubmitTab() {
      this.$refs.tabForm.validate((valid) => {
        if (valid) this.submitData('tab', this.tabForm, true)
      })
    },
    handleSubmit() {
      this.$refs.form.validate((valid) => {
        if (valid) this.submitData('platform', this.form)
      })
    },
    handelSubmitLogin() {
      this.$refs.loginForm.validate((valid) => {
        if (valid) this.submitData('login', this.loginForm)
      })
    }
  }
}
</script>

<style lang="scss" scoped>
  .docPage-card {
    margin-bottom: 20px;
    .card-title {
      font-size: 18px;
      font-weight: bold;
    }
    .card-footer {
      text-align: right;
    }
  }
  .avatar-uploader /deep/ .el-upload {
    text-align: left !important;
  }
  .avatar-uploader-icon {
    font-size: 24px;
    color: #8c939d;
    width: 120px;
    height: 120px;
    line-height: 120px;
    text-align: center;
    border: 1px dashed rgb(217, 217, 217);
    border-radius: 5px;
  }
  .avatar {
    width: 120px;
    height: 120px;
    display: inline-block;
    vertical-align: middle;
    border: 1px dashed rgb(217, 217, 217);
    border-radius: 5px;
    object-fit: cover;
    background-color: #f8f8f8;
  }
  .el-upload {
    width: 100%;
  }
</style>
