<template>
    <div>
      <el-dialog
        top="10vh"
        title="添加敏感词"
        :close-on-click-modal="false"
        :visible.sync="dialogVisible"
        width="30%"
        :before-close="handleClose"
        >
        <el-form
            :model="ruleForm"
            ref="ruleForm"
            label-width="140px"
            class="demo-ruleForm"
            :rules="rules"
            @submit.native.prevent
        >
            <el-form-item class="itemCenter">
              <el-radio-group v-model="ruleForm.importType">
                <el-radio-button :label="0">单条添加</el-radio-button>
                <el-radio-button :label="1">批量上传</el-radio-button>
              </el-radio-group>  
            </el-form-item>
            <el-form-item
            label="敏感词表名"
            prop="word"
            v-if="ruleForm.importType === 0"
            >
            <el-input
                v-model="ruleForm.word"
                placeholder="您可添加一个词"
            ></el-input>
            </el-form-item>
            <el-form-item
            label="敏感词类型"
            prop="sensitiveType"
            v-if="ruleForm.importType === 0"
            >
            <el-select v-model="ruleForm.sensitiveType" placeholder="请选择" style="width:100%;">
                <el-option
                v-for="item in sensitiveTypeOptions"
                :key="item.value"
                :label="item.name"
                :value="item.value">
                </el-option>
            </el-select>
            </el-form-item>
            <el-form-item
            label="批量上传"
            prop="fileName"
            v-if="ruleForm.importType === 1"
            >
            <el-upload
                class="upload-box"
                drag
                action=""
                :show-file-list="false"
                :auto-upload="false"
                :limit="5"
                multiple
                accept=".xlsx"
                :file-list="fileList"
                :on-change="uploadOnChange"
              >
              <div>
                <div>
                    <img :src="require('@/assets/imgs/uploadImg.png')" class="upload-img" />
                    <p class="click-text">将文件拖到此处，或<span class="clickUpload">点击上传</span></p>
                </div>
              </div>
              </el-upload>
               <!-- 上传文件的列表 -->
                <div class="file-list" v-if="fileList.length > 0">
                    <transition name="el-zoom-in-top">
                    <ul class="document_lise">
                    <li
                        v-for="(file, index) in fileList"
                        :key="index"
                        class="document_lise_item"
                    >
                        <div style="padding:8px 0;" class="lise_item_box">
                        <span class="size">
                            <img :src="require('@/assets/imgs/fileicon.png')" />
                            {{ file.name }}
                            <span class="file-size">
                            {{ filterSize(file.size) }}
                            </span>
                            <el-progress 
                            :percentage="file.percentage" 
                            v-if="file.percentage !== 100"
                            :status="file.progressStatus"
                            max="100"
                            class="progress"
                            ></el-progress>
                        </span>
                        <span class="handleBtn">
                            <span>
                            <span v-if="file.percentage === 100">
                                <i class="el-icon-check check success" v-if="file.progressStatus === 'success'"></i>
                                <i class="el-icon-close close fail" v-else></i>
                            </span>
                            <i class="el-icon-loading" v-else-if="file.percentage !== 100 && index === fileIndex"></i>
                            </span>
                            <span style="margin-left:30px;">
                            <i class="el-icon-error error" @click="handleRemove(file,index)"></i>
                            </span>
                        </span>
                        </div>
                    </li>
                    </ul>
                </transition>
                </div>
            </el-form-item>
        </el-form>
        <span
            slot="footer"
            class="dialog-footer"
        >
            <el-button 
                @click="handleClose()">
                {{$t('common.confirm.cancel')}}
            </el-button>
            <el-button
                type="primary"
                @click="submitForm('ruleForm')"
            >{{$t('common.confirm.confirm')}}</el-button>
        </span>
        </el-dialog>
    </div>
</template>
<script>
import uploadChunk from "@/mixins/uploadChunk";
import { uploadSensitiveWord } from "@/api/safety";
export default {
    mixins: [uploadChunk],
    data(){
        return{
            sensitiveTypeOptions:[
                {
                    value:'Political',
                    name:'涉政'
                },
                {
                    value:'Abuse',
                    name:'辱骂涉黄'
                },
                {
                    value:'Terror',
                    name:'暴恐'
                },
                {
                    value:'Banned',
                    name:'违禁'
                },
                {
                    value:'Security',
                    name:'信息安全'
                },
                {
                    value:'Other',
                    name:'其他'
                }
            ],
            title:"新建词表",
            dialogVisible:false,
            ruleForm:{
                importType:0,
                word:'',
                sensitiveType:'',
                fileName:'',
                tableId:''
            },
            fileList:[],
            rules: {
                word: [{ required: true, message:'请输入敏感词', trigger: "blur" }],
                sensitiveType:[{ required: true, message: '请输选择敏感词类型', trigger: "blur" }],
                fileName:[{ required: true, message: '请上传文件', trigger: "blur" }]
            }
        }
    },
    methods:{
        uploadOnChange(file, fileList){
            this.fileList = [];
            this.fileList.push(file);
            if(this.fileList.length > 0){
                this.maxSizeBytes = 0;
                this.isExpire = true;
                this.startUpload();
            }
        },
        uploadFile(chunkFileName){
            this.ruleForm.fileName = chunkFileName;
        },
        handleClose(){
            this.dialogVisible = false;
            this.clearform()
        },
        clearform(){
            this.tableId = ''
            this.$refs.ruleForm.resetFields()
            this.$refs.ruleForm.clearValidate()
        },
        submitForm(formName){
            this.$refs[formName].validate((valid) =>{
                if(valid){
                   
                }else{
                    return false;
                }
            })
        },
        showDialog(tableId){
            this.dialogVisible = true;
            this.ruleForm.tableId = tableId;
        }
    }
}
</script>
<style lang="scss" scoped>
.itemCenter{
    display:flex;
    justify-content: center;
    /deep/.el-form-item__content{
        margin-left: 0 !important;
    }
}
.upload-box{
    .upload-img{
        width:56px;
        height:56px;
        margin-top: 10px;
    }
    .clickUpload{
       color: #384bf7;
       font-weight: bold;
    }
}
.file-list{
  padding: 20px;
  .document_lise_item{
    cursor: pointer;
    padding:5px 10px;
    list-style: none;
    background: #fff;
    border-radius:4px;
    box-shadow: 1px 2px 2px #ddd;
    display:flex;
    align-items:center;
    margin-bottom:10px;
    .lise_item_box{
      width:100%;
      display:flex;
      align-items:center;
      justify-content:space-between;
      .size{
          display:flex;
          align-items:center;
          .progress{
            width:400px;
            margin-left:30px;
          }
          img{
            width: 18px;
            height:18px;
            margin-bottom:-3px;
          }
          .file-size{
            margin-left:10px;
          }
      }
    }
  }
  .document_lise_item:hover{
    background: #ECEEFE;
  }
}
</style>