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

  setTonweb = (tonweb) => {
    this.tonweb = tonweb;
  }

  setTonClient = (tonClient) => {
    this.tonClient = tonClient;
  }


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
