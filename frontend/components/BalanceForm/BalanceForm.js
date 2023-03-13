import React, {useEffect, useState} from 'react';
import Input from '../Input/Input';
import {Button} from '../Button/Button';
import {appStore} from '../../stores/appStore';
import * as styles from './BalanceForm.module.scss';
import {observer} from 'mobx-react';
import {Modal, ModalContent, ModalOverlay, useDisclosure,} from '@chakra-ui/react';
import {CloseIcon} from '../../assets';

const BalanceForm = () => {
  const {currentUser} = appStore;
  const [value, setValue] = useState('0');
  const {isOpen, onOpen, onClose} = useDisclosure();
  const [topUpBalance, setTopUpBalance] = useState('1');

  useEffect(() => {
    if (currentUser) {
      // TODO: get deposited amount
    }
  }, [currentUser]);

  const withdraw = () => {
    // TODO: withdraw tokens
  };


  const onTopUp = (e) => {
    e.preventDefault();
    // TODO: add balance
  };

  return (
    <>
      <div className={styles.container}>
        <Input
          label={'Balance'}
          value={value}
          onChange={undefined}
          disabled={!currentUser}
          inputDisabled
        />

        <div className={styles.buttonContainer}>
          <Button
            text={'Withdraw'}
            onClick={withdraw}
            disabled={currentUser?.balance <= 0 || !currentUser}
          />
          <Button
            text={'Top-up'}
            onClick={onOpen}
            disabled={!currentUser}
          />
        </div>
      </div>

      <Modal
        isCentered
        isOpen={isOpen}
        onClose={onClose}
      >
        <ModalOverlay
          bg="rgba(0,0,0,0.5)"
          backdropFilter="blur(5px)"
        />
        <ModalContent>
          <div className={styles.modalContent}>
            <div className={styles.modalContainer}>
              <button
                className={styles.closeButton}
                onClick={onClose}
              >
                <img
                  src={CloseIcon}
                  alt="Close"
                />
              </button>
              <h3>Top-up balance</h3>

              <Input
                value={topUpBalance}
                onChange={(val) => {
                  setTopUpBalance(val);
                }}
                postfix={' TON'}
              />

              <div className={styles.modalButtonContainer}>
                <Button
                  text={'Pay'}
                  onClick={onTopUp}
                  disabled={!currentUser}
                />
              </div>
            </div>
          </div>
        </ModalContent>

      </Modal>
    </>
  );
};

export default observer(BalanceForm);
