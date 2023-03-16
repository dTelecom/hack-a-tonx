import React, {useEffect, useMemo, useState} from 'react';
import Input from '../Input/Input';
import {Button} from '../Button/Button';
import {appStore} from '../../stores/appStore';
import * as styles from './NodeForm.module.scss';
import classNames from 'classnames';
import {observer} from 'mobx-react';
import {Box} from '@chakra-ui/react';
import TonWeb from "tonweb";

const readIntFromBitString = (bs, cursor, bits) => {
    let n = BigInt(0);
    for (let i = 0; i < bits; i++) {
        n *= BigInt(2);
        n += BigInt(bs.get(cursor + i));
    }
    return n;
}

const parseAddress = cell => {
    let n = readIntFromBitString(cell.bits, 3, 8);
    if (n > BigInt(127)) {
        n = n - BigInt(256);
    }
    const hashPart = readIntFromBitString(cell.bits, 3 + 8, 256);
    if (n.toString(10) + ":" + hashPart.toString(16) === '0:0') return null;
    const s = n.toString(10) + ":" + hashPart.toString(16).padStart(64, '0');
    return new TonWeb.utils.Address(s);
};

const NodeForm = () => {
  const {currentUser, tonweb} = appStore;
  const [value, setValue] = useState('');
  const [node, setNode] = useState(null);
  const [address, setAddress] = useState('');
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    if (currentUser) {
        const f = async () => {
            const master = new TonWeb.utils.Address('EQDiSqUGDaJwY3Tr5Fo8L_oMw7NaR6tJfPT4VSB-Oqw9qbwY');
            const userAddress = new TonWeb.utils.Address(currentUser.account.address);
            const cell = new TonWeb.boc.Cell();
            cell.bits.writeAddress(userAddress);
            const res = await tonweb.provider.call2(master.toString(), 'get_node_wallet_address', [['tvm.Slice', TonWeb.utils.bytesToBase64(await cell.toBoc(false))]])
            const nodeAddress = parseAddress(res);
            const nodeContractData  = await tonweb.provider.call2(nodeAddress.toString(), 'get_wallet_data', []);
        }
        
        f();

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
