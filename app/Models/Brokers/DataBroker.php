<?php namespace Models\Brokers;


class DataBroker extends Broker
{
    public function getBestRankings(): array
    {
        $sql = '
            with average as (
                select v3.route, AVG(v3.value) value from average_point v3
                join (
                    select v1.route, v2.average, 1.96 * SQRT(AVG(POW(v2.average - v1.value, 2))) as std
                    from average_point v1
                    join (
                        select route, AVG(value) average
                        from average_point group by route
                    ) v2 on v1.route = v2.route
                    group by v1.route
                ) v4 on v3.route = v4.route
                where v3.value between v4.average - v4.std and v4.average + v4.std
                group by v3.route
            )
            select * from average
            order by value asc, route asc
            limit 10';
        return $this->select($sql);
    }

    public function getWorstRankings(): array
    {
        $sql = '
            with average as (
                select v3.route, AVG(v3.value) value from average_point v3
                join (
                    select v1.route, v2.average, 1.96 * SQRT(AVG(POW(v2.average - v1.value, 2))) as std
                    from average_point v1
                    join (
                        select route, AVG(value) average
                        from average_point group by route
                    ) v2 on v1.route = v2.route
                    group by v1.route
                ) v4 on v3.route = v4.route
                where v3.value between v4.average - v4.std and v4.average + v4.std
                group by v3.route
            )
            select * from average
            order by value desc, route asc
            limit 10';
        return $this->select($sql);
    }

    public function findByPlaceId(string $placeId): ?\stdClass
    {
        $sql = '
        with average as (
            select v3.route, AVG(v3.value) value from average_point v3
            join (
                select v1.route, v2.average, 1.96 * SQRT(AVG(POW(v2.average - v1.value, 2))) as std
                from average_point v1
                join (
                    select route, AVG(value) average
                    from average_point group by route
                ) v2 on v1.route = v2.route
                group by v1.route
            ) v4 on v3.route = v4.route
            where v3.value between v4.average - v4.std and v4.average + v4.std
            group by v3.route
        ), ranking as (
            select route, row_number() over (
                order by value asc, route asc
            ) rank from average
        ), local as (
            select route, MAX(value) max, MIN(value) min
            from average_point group by route
        ), global as (
            select MAX(value) max, MIN(value) min from average
        )
        select a.route, a.place_id, r.rank, r.maxRank, if(l.min = l.max, 50, ((a.value - l.min) / (l.max - l.min)) * 100) local_value, if(g.min=g.max, 50, ((a2.value - g.min) / (g.max - g.min)) * 100) global_value from average_point a
        join global g on 1=1
        join ranking r on r.route = a.route
        join local l on l.route = a.route
        join average a2 on a2.route = a.route
        where a.place_id = ?';
        return $this->selectSingle($sql, [$placeId]);
    }

    public function findByName(string $name): ?\stdClass
    {
        $sql = '
        with average as (
            select v3.route, AVG(v3.value) value from average_point v3
            join (
                select v1.route, v2.average, 1.96 * SQRT(AVG(POW(v2.average - v1.value, 2))) as std
                from average_point v1
                join (
                    select route, AVG(value) average
                    from average_point group by route
                ) v2 on v1.route = v2.route
                group by v1.route
            ) v4 on v3.route = v4.route
            where v3.value between v4.average - v4.std and v4.average + v4.std
            group by v3.route
        ), ranking as (
            select route, row_number() over (
                order by value asc, route asc
            ) rank from average
        ), global as (
            select MAX(value) max, MIN(value) min from average
        )
        select a.route, a.place_id, r.rank, r.maxRank, if(g.min=g.max, 50, ((a2.value - g.min) / (g.max - g.min)) * 100) global_value from average_point a
        join global g on 1=1
        join ranking r on r.route = a.route
        join average a2 on a2.route = a.route
        where a.route = ?';
        return $this->selectSingle($sql, [$name]);
    }
}

/*
with average as (
    select v3.route, AVG(v3.value) value from average_point v3
    join (
        select v1.route, v2.average, 1.96 * SQRT(AVG(POW(v2.average - v1.value, 2))) as std
        from average_point v1
        join (
            select route, AVG(value) average
            from average_point group by route
        ) v2 on v1.route = v2.route
        group by v1.route
    ) v4 on v3.route = v4.route
    where v3.value between v4.average - v4.std and v4.average + v4.std
    group by v3.route
)
select route, row_number() over (
    order by value asc, route asc
) rank from average
 */