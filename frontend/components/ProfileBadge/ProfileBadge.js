import React from 'react';
import {observer} from 'mobx-react';
import {ExitIcon, WalletIcon} from '../../assets';
import {appStore} from '../../stores/appStore';
import * as styles from './ProfileBadge.module.scss';
import { CHAIN } from '@tonconnect/sdk';
// import { Address } from 'ton';

const ProfileBadge = ({signOut}) => {
  const {currentUser} = appStore;

//   const userFriendlyAddress = Address.parseRaw(currentUser.account.address).toFriendly({ testOnly: currentUser.account.chain === CHAIN.TESTNET });

//   const shortAccountId = userFriendlyAddress.slice(0, 4) + '...' + userFriendlyAddress.slice(-3);
  const shortAccountId = currentUser.account.address;

  return (
    <div className={styles.badge}>
      <div className={styles.accountBadge}>
        <img
          src={WalletIcon}
          alt={'wallet icon'}
        />
        <span>{shortAccountId}</span>
      </div>

      <button
        onClick={signOut}
        className={styles.logOutButton}
      >
        <img
          src={ExitIcon}
          alt={'log out button'}
        />
      </button>
    </div>
  );
};

export default observer(ProfileBadge);
