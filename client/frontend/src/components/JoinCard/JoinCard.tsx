import {Button} from "../Button/Button";
import styles from './JoinCard.module.scss';

interface IProps {
  img: string;
  buttonText: string;
  text: string;
  onClick: () => void;
}

export const JoinCard = (
  {
    text,
    buttonText,
    img,
    onClick
  }: IProps) => {
  return (
    <div className={styles.container}>
      <img
        src={img}
        alt={'join card screenshot'}
      />

      <p>{text}</p>

      <div className={styles.button}>
        <Button
          text={buttonText}
          onClick={onClick}
        />
      </div>
    </div>
  )
}