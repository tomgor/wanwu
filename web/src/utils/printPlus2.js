/*根据finish为1或2时，判断是否打印结束*/
import workerTimer from './worker'
import {parseSub, isSub} from "@/utils/util.js"

const Print = function (opt) {
    this.sentenceArr = []//存储待打印的句子的数组
    this.sIndexMap={} 
    this.timer = opt.timer || 10; //打印速度
    this.t = null;
    this.sIndex = 0 //记录已打印句子的索引（避免重复打印）
    this.printStatus = 0
    this.fullWord = ''
    this.searchList = []
    this.onPrintEnd = (opt.onPrintEnd && (typeof opt.onPrintEnd === 'function')) ? opt.onPrintEnd : () => {
    };
    this.looper = null
}
Print.prototype = {
    print(sentence,privateData, printingCB, endCB) {
        if(privateData.searchList  && privateData.searchList.length){
            this.searchList = privateData.searchList
        }
        this.sentenceArr.push(sentence)
        this.loop(printingCB, endCB, "truely")
    },
    stop() {
        this.sentenceArr = []
        this.sIndexMap = {}
        this.sIndex = 0
        this.looper && this.looper.stop()
    },
    loop(printingCB, endCB) {

        //如果正在打印或者打印结束
        if (this.printStatus === 1 || this.sIndex >= this.sentenceArr.length) {
            return;
        }

        let curSentence = this.sentenceArr[this.sIndex]
        this.printStatus = 1
        if(!curSentence){
            console.log(this.sIndex, this.sentenceArr)
            return;
        }
        this.looper = new Looper(this.sIndex, curSentence, this.timer, (world) => {
            this.printStatus = 1
            let isEnd = this.sIndex === this.sentenceArr.length -1
            printingCB({world,finish:curSentence.finish, isEnd},this.searchList)
        }, (data) => {
            this.printStatus = 0
            this.sIndex++;
            if (this.sentenceArr[this.sIndex]) {
                this.loop(printingCB, endCB)
            } else {
                this.onPrintEnd()
            }
        },this.sIndexMap)
    },
    getAllworld(){
        let str = ''
        this.sentenceArr.forEach((n,i)=>{
            str += n.response
        })
        return str
    }
}

const Looper = function (sIndex, sentence, timer, printCB, endCB,sIndexMap) {
    this.sIndex = sIndex
    this.sIndexMap=sIndexMap
    this.sentence = sentence ? sentence.response : "" //当前要打印的句子
    this.timer = timer
    this.t = null
    this.index = 0 //当前打印到的字符位置
    this.printCB = printCB //每打印一个字符的回调
    this.endCB = endCB //句子打印结束的回调
    this.isCodeBlock = false // 新增：标记是否为代码块
    this.codeBlockContent = '' // 新增：存储代码块内容

    // 在初始化时检测是否为代码块
    this.detectCodeBlock()
    this.start()
}

Looper.prototype = {
    detectCodeBlock() {
        // 更宽松的代码块匹配正则
        const codeBlockRegex = /^```([\s\S]*?)```$/s;
        const match = this.sentence.match(codeBlockRegex);
        
        if (match) {
            this.isCodeBlock = true;
            this.codeBlockContent = match[0]; // 整个代码块内容
            this.sentence = match[1]; // 代码块内部内容（去掉```）
        }
    },
    start() {
        if(this.sentence === ''){
            this.printCB('')
            this.stop()
            this.index++;
            return
        }

        if (this.isCodeBlock) {
            this.printCB('```' + this.sentence + '```');
            this.stop();
            return;
        }

        // 处理索引引文标签
        if(isSub(this.sentence)){
            this.printCB(parseSub(this.sentence))
            this.stop()
            this.index++;
            return
        }

        // this.printFn();

        
        const batchSize = 30; // 推荐每次输出30个字符
        const interval = 10; // 减少输出间隔时间
        this.t = workerTimer.setInterval(() => {
            if (this.index === this.sentence.length) {
                this.stop()
                return
            }
            const endIdx = Math.min(this.index + batchSize, this.sentence.length);
            const chunk = this.sentence.slice(this.index, endIdx);
            this.printCB(chunk);
            this.index = endIdx;
        }, interval,this)

        // let sentenceArr = this.sentence.split('')
        // if(sentenceArr.length>100){
        //     sentenceArr = this.sentence.split(',')
        // }
        // this.t = workerTimer.setInterval(() => {
        //     if (this.index === sentenceArr.length) {
        //         this.stop()
        //         return
        //     }
        //     this.printCB(sentenceArr[this.index])
        //     this.index++;
        // }, this.timer,this)
    },
    printFn(){
        let sentenceArr = this.sentence.split('')
        this.printCB(sentenceArr[this.index])
        this.index++;
        if(this.index !== sentenceArr.length){
            this.t = workerTimer.setTimeout(()=>{
                this.printFn()
            },this.timer,this)
        }else{
            this.stop()
        }
    },
    stop() {
        if(this.sIndexMap[`${this.sIndex}`]) {
            return;
        }
        this.sIndexMap[`${this.sIndex}`] = true;
        this.endCB({msg: 'end', index: this.sIndex});
        this.t && workerTimer.clearInterval(this.t);
        this.t = null;
    }
}

export default Print;