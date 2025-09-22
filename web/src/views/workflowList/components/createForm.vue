<template>
  <div class="workflow-list">
    <el-dialog
      :title="titleMap[type]"
      :visible.sync="dialogVisible"
      width="750"
      append-to-body
      :close-on-click-modal="false"
    >
      <el-form ref="form" :model="form" label-width="120px" :rules="rules">
        <el-form-item :label="$t('list.pluginPic') + ':'" prop="avatar">
          <el-upload
            class="avatar-uploader"
            action=""
            name="files"
            :show-file-list="false"
            :http-request="handleUploadImage"
            accept=".png,.jpg,.jpeg"
          ><!--:on-error="handleUploadError"-->
            <img
              class="upload-img"
              :src="form.avatar && form.avatar.path ? form.avatar.path : (defaultIcon || defaultLogo)"
            />
            <p class="upload-hint">
              {{this.$t('common.fileUpload.clickUploadImg')}}
            </p>
          </el-upload>
        </el-form-item>
        <el-form-item :label="$t('list.pluginName')+':'" prop="configName">
          <el-input
            :placeholder="$t('list.nameplaceholder')"
            v-model="form.configName"
            maxlength="30"
            show-word-limit
          ></el-input>
        </el-form-item>
        <!--<el-form-item label="插件英文名:" prop="configENName">
          <el-input
            type="input"
            placeholder="插件英文名将用于被大模型识别及调用，仅支持英文、数字、下划线，以英文字母开头"
            v-model="form.configENName"
          ></el-input>
        </el-form-item>-->
        <el-form-item :label="$t('list.pluginDesc')+':'" prop="configDesc">
          <el-input
            type="textarea"
            :placeholder="$t('list.descplaceholder')"
            v-model="form.configDesc"
            show-word-limit
            maxlength="600"
          ></el-input>
        </el-form-item>
        <!--v-if="type === 'create'"-->
        <el-form-item v-if="false" :label="$t('list.mapTypeLabel')+':'">
          <el-radio-group v-model="form.isStream">
            <el-radio :label="false">{{$t('list.normalMap')}}</el-radio>
            <!-- <el-radio :label="true">{{$t('list.streamMap')}}</el-radio> -->
          </el-radio-group>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">{{$t('list.cancel')}}</el-button>
        <el-button type="primary" @click="doPublish">{{$t('list.confirm')}}</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { createWorkFlow, copyExample, uploadFile } from "@/api/workflow";
import { mapGetters } from "vuex"

export default {
  props: {
    type: {
      type: String,
      default: "create",
    },
    editForm: {
      type: Object,
    },
  },
  data() {
    return {
      dialogVisible: false,
      defaultLogo: require("@/assets/imgs/bg-logo.png"),
      defaultIcon: '',
      form: {
        configName: "",
        configENName: "",
        configDesc: "",
        isStream: false,
        avatar: {
          key: '',
          path: ''
        },
      },
      titleMap: {
        edit: this.$t('list.editplugin'),
        create:this.$t('list.createplugin'),
        clone: this.$t('list.copy_Demo'),
      },
      workflowID: "",
      rules: {
        configName: [
          { required: true, message: this.$t('list.nameRules'), trigger: "change" },
          { max:30, message:this.$t('list.pluginNameRules'), trigger: "change" },
          {
            validator: (rule, value, callback) => {
              // 对其新版工作流名称规则
              if (/^[a-zA-Z][a-zA-Z0-9_]{0,63}$/.test(value)) {
                callback();
              } else {
                callback(
                  new Error(
                    this.$t('list.nameplaceholder')
                  )
                );
              }
            },
            trigger: "change",
          },
        ],
        configENName: [
          { required: false, message: this.$t('list.enNameRules'), trigger: "blur" },
          {
            validator: (rule, value, callback) => {
                if(!value){
                    callback();
                }else{
                    if (/^[a-zA-Z][a-zA-Z0-9_]*$/.test(value)) {
                        callback();
                    } else {
                        callback(
                            new Error(this.$t('list.enNameErrorRules'))
                        );
                    }
                }

            },
            trigger: "change",
          },
        ],
        configDesc: [
          { required: true, message: this.$t('list.pluginDescRules'), trigger: "blur" },
          { max: 600, message:this.$t('list.pluginLimitRules'),trigger: "blur"}
        ],
      },
    };
  },
  watch: {
    commonInfo:{
      handler(val) {
        const { defaultIcon = {} } = val.data || {}
        this.defaultIcon = defaultIcon.workflowIcon ? this.$basePath + '/user/api/' + defaultIcon.workflowIcon :  ''
      },
      deep: true
    },
  },
  computed: {
    ...mapGetters('user', ['commonInfo']),
  },
  methods: {
    getBase64(file) {
      return new Promise((resolve, reject) => {
        const fileReader = new FileReader()
        fileReader.onload = event => {
          const result = event.target ? event.target.result : ''
          if (!result || typeof result !== 'string') {
            reject('file read fail')
            return
          }
          resolve(result.replace(/^.*?,/, ''))
        }
        fileReader.readAsDataURL(file)
      })
    },
    getFileExtension(name) {
      const index = name.lastIndexOf('.');
      return name.slice(index + 1).toLowerCase();
    },
    async handleUploadImage(data) {
      if (data.file) {
        const base64 = await this.getBase64(data.file).catch(() => '')

        if (!base64) {
          this.handleUploadError()
          return
        }
        const res = await uploadFile({
          file_head: {
            file_type: this.getFileExtension(data.file.name),
            biz_type: 6,
          },
          data: base64,
        })
        const {upload_uri, upload_url} = res.data || {}
        this.form.avatar = {key: upload_uri || '', path: upload_url || ''}
      }
    },
    handleUploadError() {
      this.$message.error(this.$t('common.message.uploadError'))
    },
    openDialog(row) {
      if (this.type === "edit" && this.editForm) {
        this.form = this.editForm;
      } else {
        this.clearForm();
      }
      if(row){
        this.workflowID = row.id;
      }
      this.dialogVisible = true;
      this.$nextTick(()=>{
          this.$refs['form'].clearValidate()
      })
    },
    clearForm() {
      this.form = {
        configName: "",
        configENName: "",
        configDesc: "",
        avatar: {
          key: '',
          path: ''
        },
        isStream:false
      };
    },
    async doPublish() {
      let valid = false;
      await this.$refs.form.validate((vv) => {
        if (vv) {
          valid = true;
        }
      });
      if (!valid) return;
      if (this.type == "edit") {
        this.$emit("save");
        this.dialogVisible = false;
        return;
      }
      if (this.type == "clone") {
        this.form.workflowID = this.workflowID;
        let res = await copyExample(this.form);
        if (res.code === 0) {
          this.$message.success(this.$t('list.copySuccess'));
          this.dialogVisible = false;
          let { workflowID } = res.data;
          this.$router.push({ path: "/workflow", query: { id: workflowID } });
        }
        return;
      }
      const { configName, configDesc, avatar } = this.form || {}
      const res = await createWorkFlow({
        name: configName,
        desc: configDesc,
        avatar,
      });
      if (res.code === 0) {
        this.$message.success(this.$t('list.createSuccess'));
        this.dialogVisible = false;
        let { workflow_id } = res.data;
        let querys = { id: workflow_id };
        if(this.form.isStream){
          querys.isStream = true
        }
        this.$router.push({ path: "/workflow", query: querys });
      }
    },
  },
};
</script>

<style lang="scss" scoped>
@import "../../../style/workflow.scss";
.workflow-list {
  position: absolute;
}
.avatar-uploader {
  position: relative;
  width: 98px;
  .upload-img {
    object-fit: cover;
    width: 100%;
    height: 98px;
    background: #eee;
    border-radius: 8px;
    border: 1px solid #DCDFE6;
    display: inline-block;
    vertical-align: middle;
  }
  .upload-hint {
    position: absolute;
    width: 100%;
    bottom: 0;
    background: $color_opacity;
    color: $color;
    font-size: 12px;
    line-height: 26px;
    z-index: 10;
    border-radius: 0 0 8px 8px;
  }
}
</style>
