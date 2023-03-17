import React, {useEffect, useMemo, useState} from 'react';
import Input from '../Input/Input';
import {Button} from '../Button/Button';
import {appStore} from '../../stores/appStore';
import * as styles from './NodeForm.module.scss';
import classNames from 'classnames';
import {observer} from 'mobx-react';
import {Box} from '@chakra-ui/react';
import { MasterContract } from '../../MasterContract';
import { Address } from 'ton';

const NodeForm = () => {
  const {currentUser, tonClient} = appStore;
  const [value, setValue] = useState('');
  const [node, setNode] = useState(null);
  const [address, setAddress] = useState('');
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    if (currentUser) {
        const masterContract = tonClient.open(new MasterContract(Address.parseFriendly('EQDiSqUGDaJwY3Tr5Fo8L_oMw7NaR6tJfPT4VSB-Oqw9qbwY')));
        masterContract.getNodeWalletAddress();
    //     tonClient.getBalance(currentUser.account.address).then(balance => {
    //         setValue(balance);
    //     });
    //   // TODO: get node, stacked amount, address
    //   const node = '';
    //   if (node) {
    //     setNode(node);
    //     setValue('0');
    //     setAddress('address');
    //   } else {
    //     setValue('1');
    //   }
    }
  }, [currentUser]);


  const onSubmitNode = () => {
    if (!address) {
      return;
    }
    // TODO: add node, then request balance
  };

  const removeNode = () => {
    // TODO: remove node, then request balance
  };

  const isAdd = useMemo(() => {
    return !node;
  }, [node]);


  const onButtonClick = () => {
    if (isAdd) {
      onSubmitNode();
    } else {
      removeNode();
    }
  };

  if (node === undefined) {
    return null;
  }

  return (
    <div className={styles.container}>
      <h3>{!isAdd ? 'Current Node' : 'Adding a Node'}</h3>
      <Input
        label={isAdd ? 'Enter staking amount' : 'Staked'}
        value={value}
        onChange={!isAdd ? undefined : setValue}
        postfix={' TON'}
        disabled={!currentUser || !loaded}
      />

      <Box mt={'16px'}>
        <Input
          label={isAdd ? 'Enter address' : 'Address'}
          value={address}
          onChange={!isAdd ? undefined : setAddress}
          disabled={!currentUser || !loaded}
        />
      </Box>

      <div className={classNames(styles.buttonContainer, (!address) && styles.disabled)}>
        <Button
          text={isAdd ? 'ADD NODE' : 'DELETE NODE'}
          onClick={onButtonClick}
          disabled={!currentUser || !loaded}
        />
        <p className={classNames((!currentUser || !loaded) && styles.disabled)}>{isAdd ? 'and stake TON' : 'and take TON'}</p>
      </div>
    </div>
  );
};

export default observer(NodeForm);
