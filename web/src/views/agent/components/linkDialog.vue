<template>
    <div>
        <el-dialog
            title="联网检索配置"
            :visible.sync="dialogVisible"
            width="50%"
            :before-close="handleClose">
            <el-form ref="form" :model="form" label-width="80px">
                <el-form-item 
                label="Url"
                prop="searchUrl"
                :rules="[{ required: true, message: '请输入Url', trigger: 'blur' }]"
                >
                    <el-input v-model="form.searchUrl"></el-input>
                </el-form-item>
                <el-form-item 
                prop="searchKey"
                label="Key"
                :rules="[{ required: true, message: '请输入Key', trigger: 'blur' }]"
                >
                    <el-input v-model="form.searchKey"></el-input>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="dialogVisible = false">取 消</el-button>
                <el-button type="primary" @click="handleSubmit('form')">确定</el-button>
            </span>
        </el-dialog>
    </div>
</template>
<script>
export default {
    props:{
        linkform:{
            type:Object,
            default:null
        }
    },
    data(){
        return{
            dialogVisible:false,
            form:{
                searchUrl:'',
                searchKey:''
            }
        }
    },
    methods:{
        handleClose(){
            this.dialogVisible = false
        },
        handleSubmit(formName){
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    this.dialogVisible = false
                    this.$emit('setLinkSet',this.form)
                } else {
                    console.log('error submit!!');
                    return false;
                }
            });
            
        },
        showDialog(){
            this.dialogVisible = true
             if(this.linkform !== null){
                this.form.searchUrl = this.linkform.searchUrl;
                this.form.searchKey = this.linkform.searchKey;
            }
        }
    }
}
</script>