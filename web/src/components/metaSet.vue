<template>
    <div class="metaSet">
        <div class="tool-typ">
            <el-button icon="el-icon-plus" type="primary" @click="addMataItem" size="small">新增条件</el-button>
            <el-switch v-model="metaDataFilterParams.filterEnable" active-color="#384BF7"></el-switch>
        </div>
        <div class="docMetaData">
            <div :class="['docMetaBox',metaDataFilterParams.metaFilterParams.length > 1 ? 'docMetaContainer':'']">
                <div
                    v-for="(item,index) in metaDataFilterParams.metaFilterParams"
                    class="docItem"
                >
                    <div class="docItem_data">
                        <span class="docItem_data_label">
                            <span>Key:</span>
                        </span>
                        <el-select
                            v-model="item.key"
                            placeholder="请选择"
                            @change="keyChange($event,item)"
                        >
                            <el-option
                            v-for="item in keyOptions"
                            :key="item.key"
                            :label="item.key"
                            :value="item.key"
                            >
                            </el-option>
                        </el-select>
                    </div>
                    <el-divider direction="vertical"></el-divider>
                    <div class="docItem_data">
                        <span class="docItem_data_label">type:</span>
                        <span style="min-width:80px;">{{item.type}}</span>
                    </div>
                    <el-divider direction="vertical"></el-divider>
                    <div class="docItem_data">
                        <span class="docItem_data_label">条件:</span>
                        <el-select
                            v-model="item.condition"
                            placeholder="请选择"
                            style="width:100px;"
                        >
                            <el-option
                            v-for="item in conditionOptions[item.type]"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value"
                            >
                            </el-option>
                        </el-select>
                    </div>
                    <el-divider direction="vertical"></el-divider>
                    <div class="docItem_data">
                        <span class="docItem_data_label">value:</span>
                        <span v-if="!item.showEdit" style="min-width:120px;">{{item.metaValue}}</span>
                        <div v-else style="min-width:120px;">
                            <el-input
                                v-model="item.value"
                                v-if="item.type === 'string'"
                                @blur="metaValueBlur(item)"
                                placeholder="string"
                            ></el-input>
                            <el-input
                                v-model="item.value"
                                v-if="item.type === 'number'"
                                @blur="metaValueBlur(item)"
                                type="number"
                                placeholder="number"
                            ></el-input>
                            <el-date-picker
                                v-if="item.type === 'time'"
                                v-model="item.value"
                                align="right"
                                format="yyyy-MM-dd HH:mm:ss"
                                value-format="timestamp"
                                type="datetime"
                                placeholder="选择日期时间"
                            >
                            </el-date-picker>
                        </div>
                    </div>
                    <el-divider direction="vertical"></el-divider>
                    <div class="docItem_data docItem_data_btn">
                        <span 
                            class="el-icon-edit-outline setBtn"
                            @click="editMataItem(item)"
                            ></span>
                        <span
                            class="el-icon-delete setBtn"
                            @click="delMataItem(index)"
                        ></span>
                    </div>
                </div>
                <el-select
                    v-if="metaDataFilterParams.metaFilterParams.length > 1"
                    v-model="metaDataFilterParams.filterLogicType"
                    class="orAnd"
                    placeholder="条件"
                >
                    <el-option
                    v-for="item in conditions"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                    >
                    </el-option>
                </el-select>
            </div>
        </div>
    </div>
</template>
<script>
import {metaSelect} from "@/api/knowledge"
export default {
    props:{
        knowledgeId:{
            type:String,
            required: true,
            default:''
        }
    },
    data(){
        return {
            metaDataFilterParams:{
                filterEnable:false,
                filterLogicType:'and',
                metaFilterParams:[]
            },
            keyOptions:[],
            conditions:[
                {
                    value:'and',
                    label:'且'
                },
                {
                    value:'or',
                    label:'或'
                }
            ],
            conditionOptions:{
                time:[
                        {
                            value:'is',
                            label:'是'
                        },
                        {
                            value:'before',
                            label:'早于'
                        },
                        {
                            value:'after',
                            label:'晚于'
                        },
                        {
                            value:'empty',
                            label:'为空'
                        },
                        {
                            value:'not empty',
                            label:'不为空'
                        }
                    ],
                string:[
                    {
                        value:'is',
                        label:'是'
                    },
                    {
                        value:'is not',
                        label:'不是'
                    },
                    {
                        value:'contains',
                        label:'包含'
                    },
                    {
                        value:'not contains',
                        label:'不包含'
                    },
                    {
                        value:'start with',
                        label:'开始是'
                    },
                    {
                        value:'end with',
                        label:'结束是'
                    },
                    {
                        value:'empty',
                        label:'为空'
                    },
                    {
                        value:'not empty',
                        label:'不为空'
                    }
                ],
                number:[
                    {
                        value:'=',
                        label:'等于'
                    },
                    {
                        value:'≠',
                        label:'不等于'
                    },
                    {
                        value:'>',
                        label:'大于'
                    },
                    {
                        value:'≥',
                        label:'大于等于'
                    },
                    {
                        value:'<',
                        label:'小于'
                    },
                    {
                        value:'≤',
                        label:'小于等于'
                    },
                    {
                        value:'empty',
                        label:'为空'
                    },
                    {
                        value:'not empty',
                        label:'不为空'
                    }
                ]
            }
        }
    },
    watch:{
       'metaDataFilterParams':{
         handler: function (val) {
            if(val){
                const data = {
                    knowledgeId:this.knowledgeId,
                    metaDataFilterParams:val,
                }
                this.$emit('getMetaData',data)
            }
         },
         immediate: true,
         deep:true
       }
    },
    created(){
        this.getList()
    },
    methods:{
        getList(){
            metaSelect({knowledgeId:this.knowledgeId}).then(res =>{
                if(res.code === 0){
                    this.keyOptions = res.data.knowledgeMetaDataList || []
                }
            }).catch(() =>{})
        },
        isEmpty(value){
            if (value === null || value === undefined || value === '') return true;
            return false;
        },
        validateRequiredFields(data){
            return !data.some(obj => 
                Object.values(obj).some(val => this.isEmpty(val))
            );
        },
        editMataItem(item){
            item.showEdit = true
        },
        addMataItem(){
            if(this.metaDataFilterParams.filterEnable === false){
                this.$message.warning('请开启元数据配置后再进行添加')
                return;
            }
            if(this.metaDataFilterParams.metaFilterParams.length > 0){
                 if(!this.validateRequiredFields(this.metaDataFilterParams.metaFilterParams)){
                    this.$message.warning('存在未填信息去,请补充')
                    return
                 }
            }
            this.metaDataFilterParams.metaFilterParams.push({
                key:'',
                type:'',
                condition:'',
                value:'',
                showEdit:false
            })
        },
        clearData(){
            this.metaDataFilterParams.metaFilterParams = [];
            this.metaDataFilterParams.filterLogicType = '';
        },
        keyChange(e,item){
           item.key = e;
        },
        delMataItem(index){
            this.metaDataFilterParams.metaFilterParams.splice(index,1)
            if(this.metaDataFilterParams.metaFilterParams.length === 0){
                this.metaDataFilterParams.filterLogicType = '';
            }
        }
    }
}
</script>
<style lang="scss" scoped>
/deep/{
    .el-dialog__body{
        padding:10px 20px;
    }
}
.metaSet{
    .tool-typ{
        display: flex;
        justify-content:space-between;
        align-items:center;
    }
    .docMetaData {
        display:flex;
        justify-content:space-between;
        align-items:center;
        .docMetaContainer{
            position: relative;
            margin-left:80px;
            margin-top:15px;
        }
        .docMetaContainer::after{
            content: "";
            display: block;
            position: absolute;
            left: -20px;
            top: 50%;
            bottom: 0;
            width: 15px;
            height: 90%;
            transform: translateY(-50%);
            border-top-left-radius: 8px;
            border-bottom-left-radius: 8px;
            border: 1px solid rgba(16, 24, 40, .1411764706);
            border-right-width: 0;
        }
        .docItem {
            display: flex;
            flex:1;
            align-items: center;
            border-radius: 8px;
            background: #f7f8fa;
            margin-top: 10px;
            width: fit-content;
            .docItem_data {
            display: flex;
            align-items: center;
            flex-wrap: wrap;
            margin-bottom: 5px;
            padding: 0 10px;
            .el-input,
            .el-select,
            .el-date-picker {
                min-width: 160px;
            }
            .docItem_data_label {
                margin-right: 5px;
                display: flex;
                align-items: center;
                .question {
                color: #aaadcc;
                margin-left: 2px;
                cursor: pointer;
                }
            }
            .setBtn {
                font-size: 16px;
                cursor: pointer;
                color: #384bf7;
            }
            }
            .docItem_data_btn {
            display: flex;
            justify-content: center;
            .el-icon-delete {
                margin-left: 5px;
            }
            }
        }
        .orAnd{
            width:80px;
            position: absolute;
            left: 0;
            top: 50%;
            transform: translate(-110%, -50%);
            padding: 6px 2px;
            z-index: 1;
        }
    }
}

</style>