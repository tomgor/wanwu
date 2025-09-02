<template>
    <el-dialog
    title="新增分段"
    :visible.sync="dialogVisible"
    width="40%"
    :before-close="handleClose">
    <el-form :model="ruleForm" ref="ruleForm" label-width="100px" class="demo-ruleForm">
        <el-form-item class="itemCenter">
            <el-radio-group v-model="createType" @input="typeChage($event)">
            <el-radio-button :label="'single'">单条新增</el-radio-button>
            <el-radio-button :label="'file'">批量上传</el-radio-button>
            </el-radio-group>
        </el-form-item>
        <el-form-item label="批量上传" 
            v-if="createType === 'file'" 
            prop="fileUploadId"
            :rules="[{ required: true, message: '请上传文件', trigger: 'blur' }]"
        >
            <fileUpload ref="fileUpload" :templateUrl="templateUrl" @uploadFile="uploadFile" :accept="accept"/>
        </el-form-item>
        <template v-if="createType === 'single'">
            <el-form-item 
                label="分段内容" 
                prop="content"
                :rules="[{ required: true, message: '请输入分段内容', trigger: 'blur' }]"
                >
                <el-input placeholder="请输入分段内容" v-model="ruleForm.content" type="textarea" :rows="4"></el-input>
            </el-form-item>
            <el-form-item label="关键词" prop="labels">
                <el-tag
                :key="tag"
                v-for="(tag,index) in ruleForm.labels"
                closable
                :disable-transitions="false"
                @close="handleTagClose(index)">
                {{tag}}
                </el-tag>
                <el-input
                class="input-new-tag"
                v-if="inputVisible"
                v-model="inputValue"
                ref="saveTagInput"
                size="small"
                @keyup.enter.native="handleInputConfirm"
                @blur="handleInputConfirm"
                >
                </el-input>
                <el-button v-else class="button-new-tag" size="small" @click="showInput">+ 新建关键词</el-button>
            </el-form-item>
            <el-form-item label="新增方式">
                <el-checkbox-group v-model="checkType">
                    <el-checkbox label="more" name="type">连续新增</el-checkbox>
                </el-checkbox-group>
            </el-form-item>
        </template>
    </el-form>
    <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="submit('ruleForm')">确 定</el-button>
    </span>
    </el-dialog>
</template>
<script>
import fileUpload  from "@/components/fileUpload"
import {createSegment,createBatchSegment} from "@/api/knowledge";
export default {
    components:{fileUpload},
    data(){
        return{
            accept:'.csv',
            checkType:[],
            inputVisible:false,
            inputValue:'',
            createType:'single',
            ruleForm:{
                content:'',
                docId:'',
                labels:[],
                fileUploadId:''
            },
            dialogVisible:false,
            templateUrl:'/user/api/v1/static/docs/segment.csv'
        }
    },
    methods:{
        typeChage(val){
            if(val === 'single'){
                this.ruleForm.fileUploadId = '';
                this.$refs.fileUpload.clearFileList();
            }else{
                this.clearForm()
            }
        },
        uploadFile(fileUploadId){
            this.ruleForm.fileUploadId = fileUploadId;
        },
        handleClose(){
            this.dialogVisible = false
        },
        showDiglog(docId){
            this.dialogVisible = true
            this.ruleForm.docId = docId
        },
        showInput(){
            this.inputVisible = true;
            this.$nextTick(_ => {
                this.$refs.saveTagInput.$refs.input.focus();
            });
        },
        handleTagClose(index){
            this.ruleForm.labels.splice(index,1);
        },
        handleInputConfirm(){
            if(this.inputValue){
                this.ruleForm.labels.push(this.inputValue);
                this.inputVisible = false;
                this.inputValue = '';
            }else{
                this.$message.warning('请输入关键词')
            }
        },
        submit(formName){
            if(this.createType === 'single'){
                this.handleSingle(formName);
            }else{
                this.handleFile();
            }
        },
        handleSingle(formName){
            this.$refs[formName].validate((valid) =>{
                if(valid){
                    const data = {content:this.ruleForm.content,docId:this.ruleForm.docId,labels:this.ruleForm.labels}
                    createSegment(data).then(res =>{
                        if(res.code === 0){
                            this.$message.success('创建成功');
                            if(!this.checkType.length){
                                this.dialogVisible = false;
                            }else{
                                this.clearForm()
                            }
                            this.$emit('updateData')
                        }
                    }).catch(() =>{})
                }else{
                   return false; 
                }
            })
        },
        handleFile(){
            const data = {fileUploadId:this.ruleForm.fileUploadId,docId:this.ruleForm.docId};
            createBatchSegment(data).then(res =>{
                if(res.code === 0){
                    this.$message.success('创建成功');
                    this.dialogVisible = false;
                    this.$emit('updateDataBatch')
                }
            }).catch(() =>{})
        },
        clearForm(){
            this.ruleForm.content = ''
            this.ruleForm.labels = []
            this.checkType = []
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
.el-tag{
    margin-right: 5px;
    color:#3848f7;
    border-color: #3848f7;
    background: #f4f5ff
}
/deep/{
    .el-tag .el-tag__close{
        color: #3848f7 !important;
    }
    .el-tag .el-tag__close:hover{
        color: #fff !important;
        background: #3848f7;
    }
    .el-checkbox__input.is-checked+.el-checkbox__label{
        color: #3848f7;
    }
}
</style>