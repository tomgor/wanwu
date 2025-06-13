/*根据finish为1或2时，判断是否打印结束*/
const Print = function (opt) {
    this.sentenceArr = []
    this.timer = opt.timer || 30; //打印速度
    this.t = null;
    this.sIndex = 0
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
        this.loop(printingCB, endCB)
    },
    stop() {
        this.sentenceArr = []
        this.looper && this.looper.stop()
    },
    loop(printingCB, endCB) {

        //如果正在打印或者打印结束
        if (this.printStatus) {
            return
        }

        let curSentence = this.sentenceArr[this.sIndex]
        this.printStatus = 1
        this.looper = new Looper(this.sIndex, curSentence, this.timer, (world) => {
            this.printStatus = 1
            printingCB({world,finish:curSentence.finish},this.searchList)
        }, (data) => {
            this.printStatus = 0
            this.sIndex++;
            if (this.sentenceArr[this.sIndex]) {
                this.loop(printingCB, endCB)
            } else {
                this.onPrintEnd()
            }
        })
    },
    getAllworld(){
        let str = ''
        this.sentenceArr.forEach((n,i)=>{
            str += n.response
        })
        return str
    }
}

const Looper = function (sIndex, sentence, timer, printCB, endCB) {
    this.sIndex = sIndex
    this.sentence = sentence.response
    this.timer = timer
    this.t = null
    this.index = 0
    this.printCB = printCB
    this.endCB = endCB
    this.start()
}
Looper.prototype = {
    start() {
        if(this.sentence === ''){
            this.printCB('')
            this.stop()
            this.index++;
            return
        }
        let sentenceArr = this.sentence.split('')
        this.t = setInterval(() => {
            if (this.index === sentenceArr.length) {
                this.stop()
                return
            }
            this.printCB(sentenceArr[this.index])
            this.index++;
        }, this.timer)
    },
    stop() {
        //console.log('关闭定时器')
        this.endCB({msg: 'end', index: this.sIndex})
        this.t && clearInterval(this.t)
    }
}


export default Print
