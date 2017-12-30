
class Action {
  constructor() {
    this.handler = {}
  }

  subscribe(event, fn) {
    if (this.handler[event]) {
      this.handler[event].push(fn)
    } else {
      this.handler[event] = [fn]
    }
  }

  emit(event, ...data) {
    if (this.handler[event]) {
      for (let handle of this.handler[event]) {
        handle(...data)
      }
    }
  }
}

export default new Action()
