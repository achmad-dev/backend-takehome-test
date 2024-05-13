CREATE SCHEMA backend_test;

create table if not exists backend_test.campaigns
(
    id    SERIAL primary key,
    title varchar(255)
);

create table if not exists backend_test.payment_methods
(
    id    SERIAL primary key,
    name varchar(255)
);

create table if not exists backend_test.donations
(
    id    int not null primary key,
    payment_method_id int,
    campaign_id int,
    amount decimal(9,2)
);

insert into backend_test.campaigns (id, title) values (1, 'campaign 1');
insert into backend_test.campaigns (id, title) values (2, 'campaign 2');

insert into backend_test.payment_methods (id, name) values (1, 'payment method 1');
insert into backend_test.payment_methods (id, name) values (2, 'payment method 2');

insert into backend_test.donations(id, payment_method_id, campaign_id, amount) values (99999, 1, 1, 10000.00);