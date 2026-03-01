create table if not exists users (
    id serial primary key,
    name varchar(255) not null,
    email varchar(255) not null unique,
    age int not null default 0,
    created_at timestamptz not null default now(),
    deleted_at timestamptz
);

create table if not exists audit_logs (
    id serial primary key,
    user_id int references users(id),
    action varchar(255) not null,
    created_at timestamptz not null default now()
);

insert into users (name, email, age) values ('John Doe', 'john.doe@example.com', 28);
