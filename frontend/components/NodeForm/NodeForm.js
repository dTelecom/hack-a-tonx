import React, {useEffect, useMemo, useState} from 'react';
import Input from '../Input/Input';
import {Button} from '../Button/Button';
import {appStore} from '../../stores/appStore';
import * as styles from './NodeForm.module.scss';
import classNames from 'classnames';
import {observer} from 'mobx-react';
import {Box} from '@chakra-ui/react';

const NodeForm = () => {
  const {currentUser, contract} = appStore;
  const [value, setValue] = useState('');
  const [node, setNode] = useState(null);
  const [address, setAddress] = useState('');

  useEffect(() => {
    if (currentUser) {
      // TODO: get node, stacked amount, address
      const node = 'node';
      if (node) {
        setNode(node);
        setValue('0');
        setAddress('address');
      } else {
        setValue('10');
      }
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
        disabled={!currentUser}
      />

      <Box mt={'16px'}>
        <Input
          label={isAdd ? 'Enter IP address' : 'IP address'}
          value={address}
          onChange={!isAdd ? undefined : setAddress}
          disabled={!currentUser}
        />
      </Box>

      <div className={classNames(styles.buttonContainer, (!address) && styles.disabled)}>
        <Button
          text={isAdd ? 'ADD NODE' : 'DELETE NODE'}
          onClick={onButtonClick}
          disabled={!currentUser}
        />
        <p className={classNames(!currentUser && styles.disabled)}>{isAdd ? 'and stake TON' : 'and take TON'}</p>
      </div>
    </div>
  );
};

export default observer(NodeForm);
