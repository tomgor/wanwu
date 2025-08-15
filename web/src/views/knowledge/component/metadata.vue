<template>
  <div>
    <el-button
      icon="el-icon-plus"
      type="primary"
      size="mini"
      @click="createMetaData"
    >创建</el-button>
    <div class="docMetaData">
      <div
        v-for="(item,index) in docMetaData"
        class="docItem"
      >
        <div class="docItem_data">
          <span class="docItem_data_label">
            <span>Key:</span>
            <el-tooltip
              class="item"
              effect="dark"
              content="只能包含小写字母、数字和下划线，并且必须以小写字母开头"
              placement="top-start"
            >
              <span class="el-icon-question question"></span>
            </el-tooltip>
          </span>
          <el-input
            v-model="item.metaKey"
            @blur="metakeyBlur(item)"
          ></el-input>
        </div>
        <el-divider direction="vertical"></el-divider>
        <div class="docItem_data">
          <span class="docItem_data_label">type:</span>
          <el-select
            v-model="item.metaValueType"
            placeholder="请选择"
            @change="typeChange(item)"
          >
            <el-option
              v-for="item in typeOptions"
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
          <el-select
            v-model="item.metadataType"
            placeholder="请选择"
            style="margin-right:5px;"
            @change="typeChange(item)"
          >
            <el-option
              v-for="item in valueOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
          <el-input
            v-model="item.metaValue"
            v-if="item.metadataType ==='value' && item.metaValueType === 'string'"
            @blur="metaValueBlur(item)"
            placeholder="string"
          ></el-input>
          <el-input
            v-model="item.metaValue"
            v-if="item.metadataType ==='value'  && item.metaValueType === 'number'"
            @blur="metaValueBlur(item)"
            type="number"
            placeholder="number"
          ></el-input>
          <el-input
            v-model="item.metaRule"
            v-if="item.metadataType ==='regExp'"
            @blur="metaRuleBlur(item)"
            placeholder="regExp"
          ></el-input>
          <el-date-picker
            v-if="item.metaValueType === 'time' && item.metadataType==='value'"
            v-model="item.metaValue"
            align="right"
            format="yyyy-MM-dd HH:mm:ss"
            value-format="timestamp"
            type="datetime"
            placeholder="选择日期时间"
          >
          </el-date-picker>
        </div>
        <el-divider direction="vertical"></el-divider>
        <div class="docItem_data docItem_data_btn">
          <span
            class="el-icon-delete setBtn"
            @click="delMataItem(index)"
          ></span>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  watch: {
    docMetaData: {
      handler(val) {
        if (this.debounceTimer) {
          clearTimeout(this.debounceTimer);
        }
        this.debounceTimer = setTimeout(() => {
          this.$emit("updateMeata", val);
        }, 500);
      },
      deep: true,
      immediate: true,
    },
  },
  data() {
    return {
      debounceTimer: null,
      docMetaData: [],
      typeOptions: [
        {
          label: "String",
          value: "string",
        },
        {
          label: "Number",
          value: "number",
        },
        {
          label: "Time",
          value: "time",
        },
      ],
      valueOptions: [
        {
          value: "value",
          name: "确认值",
        },
        {
          value: "regExp",
          name: "正则表达式",
        },
      ],
    };
  },
  methods: {
    createMetaData() {
      if (this.docMetaData.length > 0 && !this.validateMetaData()) {
        return;
      }
      this.docMetaData.push({
        metaKey: "",
        metaRule: "",
        metaValue: "",
        metaValueType: "string",
        metadataType: "value",
      });
    },
    validateMetaData() {
      const hasEmptyField = this.docMetaData.some((item) => {
        const isMetaKeyEmpty =
          !item.metaKey ||
          (typeof item.metaKey === "string" && item.metaKey.trim() === "");
        const isMetaRuleRequired = item.metadataType !== "value";
        const isMetaRuleEmpty =
          isMetaRuleRequired &&
          (!item.metaRule ||
            (typeof item.metaRule === "string" && item.metaRule.trim() === ""));
        return isMetaKeyEmpty || isMetaRuleEmpty;
      });
      if (hasEmptyField) {
        this.$message.error("元数据管理存在未填写的必填字段");
        return false;
      }
      return true;
    },
    delMataItem(i) {
      this.docMetaData.splice(i, 1);
    },
    typeChange(item) {
      item.metaValue = "";
      item.metaRule = "";
    },
    metakeyBlur(item) {
      const regex = /^[a-z][a-z0-9_]*$/;
      if (!item.metaKey) {
        this.$message.warning("请输入key值");
        return;
      }
      if (!regex.test(item.metaKey)) {
        this.$message.warning("请输入符合标准的key值");
        item.metaKey = "";
        return;
      }
    },
    metaValueBlur(item) {
      if (!item.metaValue) {
        this.$message.warning("请输入value值");
        return;
      }
    },
    metaRuleBlur(item) {
      if (!item.metaRule) {
        this.showWarning("请输入正则值",item);
        return;
      }
      if (!this.isValidRegex(item.metaRule)) {
        this.showWarning("请输入合法正则值",item);
        item.metaRule = "";
        return;
      }
      if(!/^\/.*\/$/.test(item.metaRule)){
        this.showWarning("正则表达式格式应为 /pattern/ 或 /pattern/flags",item);
        item.metaRule = "";
        return;
      }
    },
    showWarning(message,item){
        this.$message.warning(message);
        item.metaRule = "";
    },
    isValidRegex(str) {
      try {
        const [_, pattern, flags] = str.match(/^\/(.*)\/([a-z]*)$/);
        new RegExp(pattern, flags);
        return true;
      } catch (e) {
        return false;
      }
    },
  },
};
</script>
<style lang="scss" scoped>
.docMetaData {
  .docItem {
    display: flex;
    align-items: center;
    border-radius: 8px;
    background: #f7f8fa;
    margin-top: 10px;
    width: fit-content;
    .docItem_data {
      display: flex;
      align-items: center;
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
}
</style>