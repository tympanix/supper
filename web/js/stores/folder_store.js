import { EventEmitter } from 'events'
import API from '../api'

class FolderStore extends EventEmitter {
  constructor() {
    super()
    this.folders = []
    this.search = ""
    this.update()
  }

  update() {
    API.getFolders().then(folders => {
      this.folders = folders
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
