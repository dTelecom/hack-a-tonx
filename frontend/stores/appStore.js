import {makeAutoObservable} from 'mobx';

class AppStore {
  currentUser = undefined;
  stats = undefined;

  constructor() {
    makeAutoObservable(this);
  }

  setCurrentUser = (currentUser) => {
    this.currentUser = currentUser;
  };


  loadStats = async () => {
    this.stats = {
      minutes: 100,
      conferences: 100,
      clients: 100,
      nodes: 100,
      income: 100,
    };
  };
}

export const appStore = new AppStore();
