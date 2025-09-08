<template>
    <div class="metaSet">
        <div class="tool-typ">
            <el-button icon="el-icon-plus" type="primary" @click="addMataItem" size="small">新增条件</el-button>
            <el-switch v-model="isEnable" active-color="#384BF7"></el-switch>
        </div>
        <div class="docMetaData">
            <div :class="['docMetaBox',docMetaData.length > 1 ? 'docMetaContainer':'']">
                <div
                    v-for="(item,index) in docMetaData"
                    class="docItem"
                >
                    <div class="docItem_data">
                        <span class="docItem_data_label">
                            <span>Key:</span>
                        </span>
                        <el-select
                            v-model="item.metaValueKey"
                            placeholder="请选择"
                            @change="keyChange($event,item)"
                        >
                            <el-option
                            v-for="item in keyOptions"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value"
                            >
                            </el-option>
                        </el-select>
                    </div>
                    <el-divider direction="vertical"></el-divider>
                    <div class="docItem_data">
                        <span class="docItem_data_label">type:</span>
                        <span style="min-width:80px;">{{item.metaValueType}}</span>
                    </div>
                    <el-divider direction="vertical"></el-divider>
                    <div class="docItem_data">
                        <span class="docItem_data_label">条件:</span>
                        <el-select
                            v-model="item.metaValueCondition"
                            placeholder="请选择"
                            style="width:100px;"
                        >
                            <el-option
                            v-for="item in conditionOptions[item.metaValueType]"
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
                                v-model="item.metaValue"
                                v-if="item.metaValueType === 'string'"
                                @blur="metaValueBlur(item)"
                                placeholder="string"
                            ></el-input>
                            <el-input
                                v-model="item.metaValue"
                                v-if="item.metaValueType === 'number'"
                                @blur="metaValueBlur(item)"
                                type="number"
                                placeholder="number"
                            ></el-input>
                            <el-date-picker
                                v-if="item.metaValueType === 'time'"
                                v-model="item.metaValue"
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
                            class="el-icon-edit-outline"
                            @click="editMataItem(item)"
                            ></span>
                        <span
                            class="el-icon-delete setBtn"
                            @click="delMataItem(index)"
                        ></span>
                    </div>
                </div>
                <el-select
                    v-if="docMetaData.length > 1"
                    v-model="conditionsKey"
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
export default {
    data(){
        return {
            isEnable:false,
            docMetaData:[],
            keyOptions:[],
            conditions:[
                {
                    value:'且',
                    label:'且'
                },
                {
                    value:'或',
                    label:'或'
                }
            ],
            conditionsKey:'',
            conditionOptions:{
                time:[
                        {
                            value:'是',
                            label:'是'
                        },
                        {
                            value:'早于',
                            label:'早于'
                        },
                        {
                            value:'晚于',
                            label:'晚于'
                        },
                        {
                            value:'为空',
                            label:'为空'
                        },
                        {
                            value:'不为空',
                            label:'不为空'
                        }
                    ],
                string:[
                    {
                        value:'是',
                        label:'是'
                    },
                    {
                        value:'不是',
                        label:'不是'
                    },
                    {
                        value:'包含',
                        label:'包含'
                    },
                    {
                        value:'不包含',
                        label:'不包含'
                    },
                    {
                        value:'开始是',
                        label:'开始是'
                    },
                    {
                        value:'结束是',
                        label:'结束是'
                    },
                    {
                        value:'为空',
                        label:'为空'
                    },
                    {
                        value:'不为空',
                        label:'不为空'
                    },
                    {
                        value:'在',
                        label:'在'
                    },
                    {
                        value:'不在',
                        label:'不在'
                    }
                ],
                number:[
                    {
                        value:'等于',
                        label:'等于'
                    },
                    {
                        value:'不等于',
                        label:'不等于'
                    },
                    {
                        value:'大于',
                        label:'大于'
                    },
                    {
                        value:'大于等于',
                        label:'大于等于'
                    },
                    {
                        value:'小于',
                        label:'小于'
                    },
                    {
                        value:'小于等于',
                        label:'小于等于'
                    },
                    {
                        value:'为空',
                        label:'为空'
                    },
                    {
                        value:'不为空',
                        label:'不为空'
                    }
                ]
            }
        }
    },
    watch:{
       docMetaData:{
         handler: function (val) {
            if(val){
                this.$emit('getMetaData',val)
            }
         },
         immediate: true,
         deep:true
       }
    },
    created(){
    },
    methods:{
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
            // if(this.docMetaData.length > 0){
            //      if(!this.validateRequiredFields(this.docMetaData)){
            //         this.$message.warning('存在未填信息去,请补充')
            //         return
            //      }
            // }
            this.docMetaData.push({
                metaValueKey:'',
                metaValueType:'',
                metaValueCondition:'',
                metaValue:'',
                showEdit:false
            })
        },
        clearData(){
            this.docMetaData = [];
            this.conditionsKey = '';
        },
        keyChange(e,item){
           item.metaValueType = e.metaValueType;
        },
        delMataItem(index){
            this.docMetaData.splice(index,1)
            if(this.docMetaData.length === 0){
                this.conditionsKey = '';
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