<template>
  <div class="login">
    <overview :login="login" />
    <div class="login-modal">
      <div class="header__left">
        <img
            v-if="login.logo && login.logo.path"
            style="height: 60px; margin: 0 15px 0 22px"
            :src="basePath + '/user/api' + login.logo.path"
        />
      </div>
<!--      <div class="container__left">-->
<!--        {{login.welcomeText}}-->
<!--      </div>-->
      <div class="login-box">
        <p class="login-header">
          <span style="font-weight: bold">{{$t('register.title')}}</span>
        </p>
        <div class="login-form">
          <el-form ref="form" :model="form" :rules="rules" label-position="top">
            <el-form-item :label="$t('register.form.username')" class="login-form-item" prop="username">
              <img class="login-icon" src="@/assets/imgs/user.png" alt="" />
              <el-input
                  v-model.trim="form.username"
                  :placeholder="$t('common.input.placeholder') + $t('register.form.username')" clearable
              />
            </el-form-item>
            <el-form-item :label="$t('register.form.email')" class="login-form-item" prop="email">
              <img class="login-icon" src="@/assets/imgs/user.png" alt="" />
              <el-input
                  v-model.trim="form.email"
                  :placeholder="$t('common.input.placeholder') + $t('register.form.email')" clearable
              />
            </el-form-item>
            <el-form-item :label="$t('register.form.code')" class="login-form-item" prop="code">
              <img class="login-icon" src="@/assets/imgs/code.png" alt="" />
              <el-input
                  style="width: calc(100% - 90px)"
                  v-model.trim="form.code"
                  @keyup.enter.native="addByEnterKey"
                  :placeholder="$t('common.input.placeholder') + $t('register.form.code')" clearable
              />
              <el-button
                  style="height: 32px; width: 80px; margin-left: 10px; vertical-align: middle; padding-left: 8px; padding-top: 8px"
                  @click="requestEmailCode({email: form.email})"
                  :disabled="isCooldown"
              >
                {{ isCooldown ? `${cooldownTime}s` : $t('register.action') + $t('register.form.code') }}
              </el-button>
            </el-form-item>
          </el-form>
          <div class="register-bt">
            <p class="primary-bt" :style="`background: ${loginButtonColor} !important`" @click="doRegister">
              {{$t('register.button')}}
            </p>
          </div>
          <div class="bottom-text">{{login.platformDesc}}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import overview from '@/views/layout/overview'
import { requestEmailCode, getCommonInfo, register } from "@/api/user"
import {replaceTitle, replaceIcon, redirectUrl} from "@/utils/util"
import router from "@/router";
export default {
  components: { overview },
  data(){
    return{
      form:{
        username:'',
        email:'',
        code: '',
      },
      rules: {
        username: [
          { required: true, message: this.$t('common.input.placeholder'), trigger: 'blur' },
          { min: 2, max: 20, message: this.$t('common.hint.userNameLimit'), trigger: 'blur'},
          { pattern: /^(?!_)[a-zA-Z0-9_.\u4e00-\u9fa5]+$/, message: this.$t('common.hint.userName'), trigger: "blur"} // 结尾：(?!.*?_$)
        ],
        email: [
          { required: true, message: this.$t('common.input.placeholder'), trigger: 'blur' },
          { pattern: /^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(.[a-zA-Z0-9_-]+)+$/, message: this.$t('common.hint.emailError'), trigger: "blur"}
        ],
        code: [
          { required: true, message: this.$t('common.input.placeholder'), trigger: 'blur' }
        ]
      },
      home: {},
      login: {},
      isCooldown: false,
      cooldownTime: 60,
      cooldownTimer: '',
      loginButtonColor: '',
      basePath: this.$basePath
    }
  },
  created() {
    // // 如果已登录，重定向到有权限的页面
    // if (this.$store.state.user.token && localStorage.getItem("access_cert")) redirectUrl()

    this.getLogoInfo()
  },
  methods:{
    addByEnterKey(e){
      if (e.keyCode === 13) {
        this.doRegister()
      }
    },
    getLogoInfo() {
      getCommonInfo().then(res => {
        const {login, home, tab = {}} = res.data || {}
        this.login = login || {}
        this.home = home || {}
        this.loginButtonColor = this.login.loginButtonColor || '#384BF7'
        replaceTitle(tab.title)
        replaceIcon(tab.logo ? tab.logo.path : '')
      })
    },
    doRegister() {
      this.$refs.form.validate((valid) => {
        if (!valid) return
        register(this.form).then(res => {
          if (res.code === 0) {
            this.$router.push({path: `/login`})
          }
        })
      })
    },
    requestEmailCode(data) {
      this.$refs.form.validateField(['email'], err => {
        if (err) return
        this.isCooldown = true
        this.cooldownTimer = setInterval(() => {
          if (this.cooldownTime > 1) {
            this.cooldownTime--
          } else {
            this.isCooldown = false
            this.cooldownTime = 60
            clearInterval(this.cooldownTimer)
          }
        }, 1000)
        requestEmailCode(data)
      })
    }
  },
  beforeDestroy() {
    clearInterval(this.cooldownTimer)
  }
}
</script>

<style lang="scss" scoped>
.login{
  height: 100%;
}
.login-modal{
  position: fixed;
  top:0;
  bottom:0;
  left:0;
  right: 0;
  width: 100%;
  height: 100%;
  z-index: 1000;
  .header__left {
    position: relative;
    width: 100%;
    min-width: 500px;
    color: #fff;
    font-weight: bold;
    display: flex;
    align-items: center;
    margin-top: 16px;
    margin-left: 10px;
    height: 60px;
  }
  .container__left {
    display: flex;
    align-items: center;
    height: calc(80% - 60px);
    font-size: 35px;
    width: calc(100% - 13% - 400px);
    justify-content: center;
    color: #fff;
    text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.6);
  }
  .login-box{
    position: absolute;
    width: 400px;
    min-width: 400px;
    height: 460px;
    top:0;
    bottom:0;
    right: 13%;
    margin: auto;
    background: rgba(244,247,255,0.5);
    border-radius: 4px;
    //box-shadow: 0px 5px 30px 0px rgba(5,8,27,0.1);
    //border: 1px solid;
    //border-image: linear-gradient(180deg, rgba(255, 255, 255, 0.3), rgba(255, 255, 255, 0.8)) 1 1;
    backdrop-filter: blur(10px);
    .login-header{
      color: $color_title;
      text-align: left;
      padding: 30px 30px 0 30px;
      //border-bottom: 1px solid rgba(235,235,235,0.8);
      span {
        display: inline-block;
        vertical-align: bottom;
        font-size: 24px;
      }
    }
    .login-form{
      padding: 30px;
      /deep/ {
        .el-input {
          height: 30px;
          line-height: 36px;
        }
        .el-input__prefix{
          left:14px;
        }
        .el-form-item__label{
          font-size: 12px!important;
          padding: 0!important;
          margin-left: 0!important;
          line-height: 26px!important;
          color: #425466;
        }
        .el-input__inner{
          font-size: 12px !important;
          background: #fff;
          border: none !important;
          padding-left: 50px;
        }
        i {
          color: #666!important;
        }
        .el-icon-lock, .el-icon-user, .el-icon-key {
          font-size: 18px;
        }
        input[type=input]::placeholder{
          color: #B3B1BC;
        }
        input::-webkit-input-placeholder{
          color: #B3B1BC;
        }
        input::-moz-placeholder{   /* Mozilla Firefox 19+ */
          color: #B3B1BC;
        }
        input:-moz-placeholder{    /* Mozilla Firefox 4 to 18 */
          color: #B3B1BC;
        }
        input:-ms-input-placeholder{  /* Internet Explorer 10-11 */
          color: #B3B1BC;
        }
      }
      .login-form-item {
        position: relative;
      }
      .login-icon {
        position: absolute;
        width: 17px;
        z-index: 10;
        top: 13.5px;
        left: 17px;
      }
      .pwd-icon {
        position: absolute;
        width: 17px;
        z-index: 10;
        top: 13.5px;
        right: 17px;
        cursor: pointer;
      }
      .register-bt{
        width: 100%;
        height: 40px;
        line-height: 30px;
        font-size: 14px;
        margin-top: 40px;
        .primary-bt {
          padding: 13px 0;
          height: 40px;
          //background: $color;
          border: none;
          font-size: 14px;
          border-radius: 3px !important;
        }
        .disabled.primary-bt {
          opacity: 0.5;
          cursor: not-allowed;
          border-color: #e9e9eb;
        }
        .operation{
          /*color: #fff;*/
          color: #8d8d8d;
          display: flex;
          margin: 30px 70px;
          span{
            display: block;
            flex:1;
            text-align: center;
            font-size: 12px;
            cursor: pointer;
          }
        }
      }
      .bottom-text {
        text-align: center;
        font-weight: normal;
        color: #888;
        margin-top: 6px;
        font-size: 12px;
      }
    }
    .login-pwd-input /deep/ {
      .el-input__inner {
        padding-right: 46px;
      }
    }
  }
}
</style>
