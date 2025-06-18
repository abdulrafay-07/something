import { format } from "date-fns";

export function getDateText(date: string) {
  const createdAt = new Date(date);
  const now = new Date();

  let dateText = "";
  const isToday =
    createdAt.getFullYear() === now.getFullYear() &&
    createdAt.getMonth() === now.getMonth() &&
    createdAt.getDate() === now.getDate();
  
  if (isToday) {
    dateText = format(createdAt, "MMMM d");
  } else if (createdAt.getFullYear() === now.getFullYear()) {
    dateText = format(createdAt, "MMMM d");
  } else {
    dateText = format(createdAt, "MMMM d, yyyy");
  }

  return dateText;
}

export function getRandomFlower() {
  const num = Math.floor(Math.random() * 3) + 1;

  return `flower${num}`;
}
