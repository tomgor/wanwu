<template>
  <el-dialog
    :visible.sync="dialogVisible"
    title="命中分段详情"
    width="70%"
    @close="handleClose"
    class="section-dialog"
  >
    <div class="section-show-container">
      <!-- 父分段区域 -->
      <div class="parent-segment" v-if="parentSegment">
        <div class="segment-header">
          <span class="parent-badge">父分段</span>
          <div class="parent-score">
            <span class="score-label">命中得分:</span>
            <span class="score-value">{{ formatScore(parentSegment.score) }}</span>
          </div>
        </div>
        <div class="parent-content">
          {{ parentSegment.content }}
        </div>
      </div>

      <!-- 子分段区域 -->
      <div class="sub-segments" v-if="segmentList.length > 0">
        <div class="segment-header">
          <span class="sub-badge">命中{{ segmentList.length }}个子分段</span>
        </div>
        <el-collapse 
          v-model="activeNames" 
          class="section-collapse"
          accordion
        >
          <el-collapse-item 
            v-for="(segment, index) in segmentList" 
            :key="index"
            :name="index"
            class="segment-collapse-item"
          >
            <template slot="title">
              <span class="segment-badge">C#-{{ index + 1 }}</span>
              <span class="segment-score">
                <span class="score-label">命中得分:</span>
                <span class="score-value">{{ formatScore(segment.score) }}</span>
              </span>
            </template>
            {{ index + 1 }}、{{ segment.content }}
            <span class="segment-action">(展示完整分段内容)</span>
            <span v-if="segment.autoSave" class="auto-save">--失去焦点自动保存</span>
          </el-collapse-item>
        </el-collapse>
      </div>
    </div>
  </el-dialog>
</template>

<script>
export default {
  name: 'SectionShow',
  data() {
    return {
      dialogVisible: false,
      activeNames: [],
      parentSegment: {
        
      },
      segmentList: [
       
      ]
    }
  },
  methods: {
    formatScore(score) {
      // 格式化得分，保留5位小数
      if (typeof score !== 'number') {
        return '0.00000';
      }
      return score.toFixed(5);
    },
    
    
    // 显示弹框
    showDiaglog(data) {
      if (data) {
        // 更新父分段数据
        if (data.searchList) {
          this.parentSegment = {
            score: parseFloat(data.score) || 0,
            content: data.searchList.snippet||'暂无内容'
          };
        }
        
        // 更新子分段数据
        if (data.searchList && Array.isArray(data.searchList.childContentList)) {
          this.segmentList = data.searchList.childContentList.map(segment => ({
            content: segment.childsnippet || '',
            autoSave: Boolean(segment.autoSave),
            score: parseFloat(segment.score) || 0
          }));
        }
        
        
      }
      this.dialogVisible = true;
    },
    // 关闭弹框
    handleClose() {
      this.dialogVisible = false;
    },
  }
}
</script>

<style lang="scss" scoped>
.section-dialog {
  /deep/ .el-dialog__body {
    padding: 0 20px 20px 20px;
    max-height:70vh;
    overflow-y: auto;
  }
}

.section-show-container {
  .parent-segment {
    padding: 20px 20px 0 20px;
    background: #fff;
    border-radius: 8px;
    
    .segment-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 15px;
      
      .parent-badge {
        background-color: #d2d7ff;
        color: #384BF7;
        padding: 6px 12px;
        border-radius: 6px;
        font-size: 12px;
        font-weight: 500;
      }
      
      .parent-score {
        display: flex;
        align-items: center;
        
        .score-label {
          font-size: 12px;
          color: #384BF7;
          font-weight: bold;
          margin-right: 5px;
        }
        
        .score-value {
          font-size: 14px;
          color: #384BF7;
          font-weight: bold;
          font-family: 'Courier New', monospace;
        }
      }
    }
    
    .parent-content {
      text-align: left;
      background-color: #f7f8fa;
      padding: 15px;
      border-radius: 6px;
      border: 1px solid #384BF7;
      
      .parent-item {
        margin-bottom: 10px;
        font-size: 14px;
        color: #333;
        line-height: 1.5;
        text-align: left;
        
        .segment-action {
          color: #999;
          font-size: 12px;
          margin-left: 8px;
        }
      }
    }
  }
  
  .sub-segments {
    padding: 20px;
    
    .segment-header {
      margin-bottom: 15px;
      
      .sub-badge {
        background-color: #d2d7ff;
        color: #384BF7;
        padding: 6px 12px;
        border-radius: 6px;
        font-size: 12px;
        font-weight: 500;
      }
    }
    
    .section-collapse {
      background-color: #f7f8fa;
      border-radius: 6px;
      border: 1px solid #384BF7;
      overflow: hidden;
      
      /deep/ .el-collapse {
        border: none;
        border-radius: 6px;
      }
      
      /deep/ .el-collapse-item__header {
        background-color: #f7f8fa;
        border-bottom: 1px solid #e4e7ed;
        padding: 12px 20px;
        font-weight: normal;
        border-left: none;
        border-right: none;
        border-top: none;
        
        &:hover {
          background-color: #f0f2f5;
        }
      }
      
      /deep/ .el-collapse-item__content {
        padding: 15px 20px;
        background-color: #fff;
        border-bottom: 1px solid #e4e7ed;
        border-left: none;
        border-right: none;
        border-top: none;
      }
      
      /deep/ .el-collapse-item:last-child .el-collapse-item__content {
        border-bottom: none;
      }
      
      /deep/ .el-collapse-item__header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        width: 100%;
        padding: 12px 20px;
        position: relative;
      }
      
      /deep/ .el-collapse-item__arrow {
        display: none !important;
      }
      
      .segment-badge {
        // background-color: #eee;
        color: #384BF7;
        // padding: 6px 12px;
        // border-radius: 6px;
        font-size: 12px;
        min-width: 40px;
        text-align: center;
        font-weight: 500;
        margin-right: 120px; // 为右边的得分留出空间
      }
      
      .segment-score {
        display: flex;
        align-items: center;
        position: absolute;
        right: 20px;
        top: 50%;
        transform: translateY(-50%);
        
        .score-label {
          font-size: 12px;
          color: #384BF7;
          font-weight: bold;
          margin-right: 5px;
        }
        
        .score-value {
          font-size: 14px;
          color: #384BF7;
          font-weight: bold;
          font-family: 'Courier New', monospace;
        }
      }
      
      /deep/ .el-collapse-item__content {
        font-size: 14px;
        color: #333;
        line-height: 1.5;
        text-align: left;
        word-wrap: break-word;
        word-break: break-all;
        overflow-wrap: break-word;
        
        .segment-action {
          color: #999;
          font-size: 12px;
          margin-left: 8px;
        }
        
        .auto-save {
          color: #666;
          font-size: 12px;
          margin-left: 8px;
          font-style: italic;
        }
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .section-show-container {
    .parent-segment,
    .sub-segments {
      padding: 15px;
    }
    
    .section-collapse {
      /deep/ .el-collapse-item__header {
        flex-direction: column;
        align-items: flex-start;
        
        .segment-badge {
          margin-bottom: 8px;
          margin-right: 0;
        }
        
        .segment-score {
          position: static;
          transform: none;
          margin-top: 8px;
        }
      }
    }
  }
}
</style>