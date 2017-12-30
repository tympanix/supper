import { EventEmitter } from 'events'
import axios from 'axios'

class FolderStore extends EventEmitter {
  constructor() {
    super()
    this.folders = []
    this.search = ""
    this.update()
  }

  update() {
    axios.get("./api/media").then(res => {
      this.folders = res.data
      this.emit("change")
    })
  }

  getAll() {
    return this.folders
  }

  getSearch() {
    return this.search
  }

  setSearch(search) {
    this.search = search
  }
}


const folderStore = new FolderStore

export default folderStore
