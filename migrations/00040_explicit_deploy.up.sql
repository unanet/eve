alter table job
    add explicit_deploy bool default false not null;

alter table service
    add explicit_deploy bool default false not null;

alter table namespace
    rename column explicit_deploy_only TO explicit_deploy;