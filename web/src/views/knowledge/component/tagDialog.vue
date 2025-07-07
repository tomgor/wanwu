<template>
    <el-dialog
        title="创建标签"
        :visible.sync="dialogVisible"
        width="20%"
        :before-close="handleClose">
        <div>
            <el-input placeholder="搜索标签" suffix-icon="el-icon-search" @keyup.enter.native="addByEnterKey"></el-input>
            <div class="add"  @click="addTag"><span class="el-icon-plus add-icon"></span>创建标签</div>
            <div class="tag-box">
                <div 
                v-for="item in tagList" 
                class="tag_item"  
                @mouseenter="mouseEnter(item)"
                @mouseleave="mouseLeave(item)"
                @dblclick="handleDoubleClick(item)"
                >
                    <el-checkbox v-model="item.checked" v-if="!item.showIpt">{{item.tagName}}</el-checkbox>
                    <el-input v-model="item.tagName" v-if="item.showIpt" @blur="inputBlur(item)" ></el-input>
                    <span class="el-icon-close del-icon" v-if="item.showDel && !item.showIpt"></span>
                </div>
            </div>
        </div>
        <span slot="footer" class="dialog-footer">
            <el-button @click="dialogVisible = false">取 消</el-button>
            <el-button type="primary" @click="dialogVisible = false">确 定</el-button>
        </span>
    </el-dialog>
</template>
<script>
export default {
    data(){
        return{
            dialogVisible:false,
            tagList:[
                {
                    tagName:'标签1',
                    checked:false,
                    showDel:false,
                    showIpt:false
                },
                {
                    tagName:'标签2',
                    checked:false,
                    showDel:false,
                    showIpt:false
                },
                {
                    tagName:'标签3',
                    checked:false,
                    showDel:false,
                    showIpt:false
                }
            ]
        }
    },
    methods:{
        showDiaglog(){
            this.dialogVisible = true;
        },
        handleClose(){
            this.dialogVisible = false;
        },
        mouseEnter(n){
            n.showDel = true
        },
        mouseLeave(n){
            n.showDel = false
        },
        handleDoubleClick(n){
            n.showIpt = true
        },
        inputBlur(n){
            n.showIpt = false
        },
        addTag(){
            this.tagList.unshift(
                {
                    tagName:'',
                    checked:false,
                    showDel:false,
                    showIpt:true
                }
            )
        },
        addByEnterKey(e){
            if (e.keyCode === 13){
                console.log('搜索标签')
            }
        }
    }
}
</script>
<style lang="scss" scoped>
/deep/{
    .el-dialog__body{
        padding:5px 20px!important;
    }
    .add{
        margin-top:10px;
        padding:10px 0;
        cursor: pointer;
      .add-icon{
        margin-right:5px;
      }
    }
    .tag-box{
        max-height:300px;
        overflow-y:scroll;
    }
    .tag_item{
        cursor: pointer;
        background:#F4F5FF;
        padding:5px;
        margin:10px 0;
        border-radius:4px;
        display:flex;
        justify-content:space-between;
        align-items:center;
        .del-icon{
            color: #384BF7;
            cursor: pointer;
            font-size: 16px;
        }
    }
}
</style>