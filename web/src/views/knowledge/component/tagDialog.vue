<template>
  <el-dialog
    title="创建标签"
    :visible.sync="dialogVisible"
    width="20%"
    :before-close="handleClose"
  >
    <div>
      <el-input
        v-if="type === 'knowledge'"
        placeholder="搜索标签"
        suffix-icon="el-icon-search"
        @keyup.enter.native="addByEnterKey"
        v-model="tagName"
      ></el-input>
      <div
        class="add"
        @click="addTag"
      ><span class="el-icon-plus add-icon"></span>创建标签</div>
      <div class="tag-box">
        <div
          v-for="(item,index) in tagList"
          class="tag_item"
          @mouseenter="mouseEnter(item)"
          @mouseleave="mouseLeave(item)"
          @dblclick="handleDoubleClick(item)"
        >
          <el-checkbox
            v-model="item.selected"
            v-if="!item.showIpt && type === 'knowledge'"
          >{{item.tagName}}</el-checkbox>
          <span v-if="!item.showIpt && type === 'doc'">{{item.tagName}}</span>
          <el-input
            v-model="item.tagName"
            v-if="item.showIpt"
            @keydown.backspace.native="handleDelete(item,index)" 
            @keyup.enter.native="inputBlur(item)"
            @blur="inputBlur(item)"
          ></el-input>
          <span
            class="el-icon-close del-icon"
            v-if="item.showDel && !item.showIpt"
            @click="type=== 'doc' ? delDocTag(index) : delTag(item,index)"
          ></span>
        </div>
      </div>
    </div>
    <span
      slot="footer"
      class="dialog-footer"
    >
      <el-button
        type="primary"
        @click="submitDialog"
      >确 定</el-button>
    </span>
  </el-dialog>
</template>
<script>
import { delTag, tagList, createTag, editTag, bindTag,bindTagCount,updateDocTag} from "@/api/knowledge";
export default {
  props:['type','tagList'],
  data() {
    return {
      dialogVisible: false,
      tagList: [],
      tagName: "",
      knowledgeId: "",
      docId:""
    };
  },
  watch:{
    tagList:{
      handler(val){
        if(val && val.length){
          this.tagList = val.map(item =>({
            ...item,
            tagName:item,
            checked: false,
            showDel: false,
            showIpt: true,
          }))
        }
      },
      deep:true,
      immediate:true
    }
  },
  methods: {
    submitDialog() {
      if(this.type === 'doc'){
        this.updateTag()
      }else{
        this.bindTag()
      }
    },
    updateTag(){
      const docTagList = this.tagList.map((item) => item.tagName)
      updateDocTag({docId:this.docId,docTagList}).then(res =>{
        if(res === 0){
          this.dialogVisible = false;
          this.$emit("relodaData");
        }
      })
    },
    bindTag(){
      const ids = this.tagList.filter((item) => item.selected).map((item) => item.tagId);
      bindTag({ knowledgeId: this.knowledgeId, tagIdList: ids }).then((res) => {
        if (res.code === 0) {
          this.$emit("relodaData");
        }
      });
      this.dialogVisible = false;
    },
    getList() {
      tagList({ knowledgeId: this.knowledgeId, tagName: this.tagName }).then(
        (res) => {
          if (res.code === 0) {
            this.tagList = res.data.knowledgeTagList.map((item) => ({
              ...item,
              showDel: false,
              showIpt: false,
            }));
          }
        }
      );
    },
    bindTagCount(tagId){
      return bindTagCount({tagId}).then(res =>{
        if(res.code === 0){
          const tagBindCount = res.data.tagBindCount;
          return tagBindCount > 0;
        }
        return 'unknow'
      }).catch(() =>{
        return 'unknow'
      })
    },
    async delTag(item) {
      const isBind = await this.bindTagCount(item.tagId);
      if(isBind == 'unknow') return
      await this.$confirm(
        `删除标签${item.tagName}`,
        item.selected && isBind ? "标签正在使用中，是否删除？" : "确认要删除当前标签？",
        {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        }
      )
        .then(async() => {
          const res = await delTag({ tagId: item.tagId })
            if (res.code === 0) {
                this.getList();
            }
        })
        .catch((error) => {
            this.getList();
        });
    },
    showDiaglog(id) {
      this.dialogVisible = true;
      if (this.type === 'doc') {
        this.docId = id;
      } else {
        this.knowledgeId = id;
        this.getList();
      }
    },
    handleClose() {
      this.dialogVisible = false;
    },
    mouseEnter(n) {
      n.showDel = true;
    },
    mouseLeave(n) {
      n.showDel = false;
    },
    handleDoubleClick(n) {
      n.showIpt = true;
    },
    inputBlur(n) {
      if(!n.tagName) return;
      if(this.type === 'doc') {
        n.showIpt = false;
        return;
      }
      if (n.tagId) {
        this.edit_tag(n);
      } else {
        this.add_Tag(n);
      }
    },
    handleDelete(n,i){
      if(n.tagName === '' && !n.tagId){
          this.tagList.splice(i,1)
      }
    },
    delDocTag(i){
      this.tagList.splice(i,1)
    },
    add_Tag(n) {
      createTag({ tagName: n.tagName }).then((res) => {
        if (res.code === 0) {
          n.showIpt = false;
          this.getList();
        }
      });
    },
    edit_tag(n) {
      editTag({ tagId: n.tagId, tagName: n.tagName }).then((res) => {
        if (res.code === 0) {
          n.showIpt = false;
          this.getList();
        }
      });
    },
    addTag() {
      const emptyTag = this.tagList.find(tag => !tag.tagId && tag.tagName === "");
      if(emptyTag) return;
      this.tagList.unshift({
        tagName: "",
        checked: false,
        showDel: false,
        showIpt: true,
      });
    },
    addByEnterKey(e) {
      if (e.keyCode === 13) {
        // 当type为doc时，不应该调用getList方法
        if (this.type === 'doc') {
          return;
        }
        this.getList();
      }
    },
  },
};
</script>
<style lang="scss" scoped>
/deep/ {
  .el-dialog__body {
    padding: 5px 20px !important;
  }
  .add {
    margin-top: 10px;
    padding: 10px 0;
    cursor: pointer;
    .add-icon {
      margin-right: 5px;
    }
  }
  .tag-box {
    max-height: 300px;
    overflow-y: scroll;
  }
  .tag_item {
    cursor: pointer;
    background: #f4f5ff;
    padding: 5px;
    margin: 10px 0;
    border-radius: 4px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    .del-icon {
      color: #384bf7;
      cursor: pointer;
      font-size: 16px;
    }
  }
}
</style>