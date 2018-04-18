import { EventEmitter } from 'events'
import API from '../api'

class FolderStore extends EventEmitter {
  constructor() {
    super()
    this.folders = []
    this.search = ""
    this.loading = true
    this.update()
  }

  update() {
    API.getFolders().then(folders => {
      this.folders = folders
      this.loading = false
      this.emit("change")
    })
  }

  getAll() {
    return this.folders
  }

  isLoading() {
    return this.loading
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
