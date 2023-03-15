import React, {useState} from 'react';
import {Header} from '../../components/Header/Header';
import {useLocation, useNavigate, useParams} from 'react-router-dom';
import ParticipantsBadge from '../../components/ParticipantsBadge/ParticipantsBadge';
import {Flex} from '@chakra-ui/react';
import styles from './JoinViewer.module.scss';
import Footer from '../../components/Footer/Footer';
import {Container} from '../../components/Container/Container';
import Input from '../../components/Input/Input';
import {FaceIcon} from '../../assets';
import {useBreakpoints} from '../../hooks/useBreakpoints';
import {Button} from '../../components/Button/Button';

export const JoinViewer = () => {
  const {isMobile} = useBreakpoints();
  const navigate = useNavigate();
  const location = useLocation();
  const {sid} = useParams();
  const [room] = useState(location.state?.room);
  const [name, setName] = useState('');

  if (!room) {
    navigate('/');
  }

  const onJoin = () => {
    navigate('/call/' + sid, {
      state: {
        name, noPublish: true, e2ee: room.e2ee,
        title: room.title,
      }
    });

  };

  return (
    <>
      <Header
        centered
        title={room?.title}
      >
        <Flex
          h={40}
          alignItems={'center'}
        >
          <span className={styles.smallText}>at the room:</span>&nbsp;<ParticipantsBadge count={room?.count}/>
        </Flex>
      </Header>

      <Container>
        <Flex
          mt={80}
          flexDirection={'column'}
          alignItems={'center'}
        >
          <h1 className={styles.title}>
            {room.hostName}{(isMobile ? '\n' : ' ') + 'invites you'}
          </h1>

          <Flex
            my={isMobile ? 24 : 40}
            gap={'30px'}
            width={'100%'}
            maxWidth={'420px'}
            alignItems={'stretch'}
          >
            <Input
              value={name}
              onChange={setName}
              label={'Enter your name'}
              icon={FaceIcon}
              placeholder={'John'}
              containerStyle={{
                width: '100%',
                textAlign: 'center',
              }}
            />
          </Flex>

          <Flex mb={12}>
            <span className={styles.subtitle}>Join as a viewer:</span>
          </Flex>

          <div className={styles.button}>
            <Button
              onClick={onJoin}
              text={'Free'}
              disabled={!name}
            />
          </div>
          {room.e2ee && (
            <p className={styles.smallGreyText}>
              End-to-end encryption is enabled,<br/>
              we recommend using the Google Chrome browser.
            </p>
          )}
        </Flex>
      </Container>

      <Footer/>
    </>
  );
};
