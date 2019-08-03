create table users(
  id serial primary key,
  username varchar(50) not null,
  email varchar(200),
  password char(60),
  created_at timestamptz default current_timestamp,
  updated_at timestamptz,
  deleted_at timestamptz
);

create table events(
  id serial primary key,
  user_id int not null,
  foreign key (user_id) references users (id),
  name varchar(200) not null,
  description varchar(500),
  scheduled timestamptz not null,
  created_at timestamptz default current_timestamp,
  updated_at timestamptz,
  deleted_at timestamptz
);
