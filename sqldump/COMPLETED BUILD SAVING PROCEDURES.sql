-- MariaDB dump 10.19  Distrib 10.11.4-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: UDR
-- ------------------------------------------------------
-- Server version	10.11.4-MariaDB-1~deb12u1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `buildfeats`
--

DROP TABLE IF EXISTS `buildfeats`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `buildfeats` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `build_id` int(10) unsigned DEFAULT NULL,
  `feat` varchar(2000) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `buildfeats_FK` (`build_id`),
  CONSTRAINT `buildfeats_FK` FOREIGN KEY (`build_id`) REFERENCES `builds` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `buildfeats`
--

LOCK TABLES `buildfeats` WRITE;
/*!40000 ALTER TABLE `buildfeats` DISABLE KEYS */;
/*!40000 ALTER TABLE `buildfeats` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `builds`
--

DROP TABLE IF EXISTS `builds`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `builds` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `client_id` int(10) unsigned DEFAULT NULL,
  `title` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `title_un` (`title`),
  KEY `builds_FK` (`client_id`),
  CONSTRAINT `builds_FK` FOREIGN KEY (`client_id`) REFERENCES `client` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=latin1 COLLATE=latin1_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `builds`
--

LOCK TABLES `builds` WRITE;
/*!40000 ALTER TABLE `builds` DISABLE KEYS */;
/*!40000 ALTER TABLE `builds` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `buildskills`
--

DROP TABLE IF EXISTS `buildskills`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `buildskills` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `build_id` int(10) unsigned DEFAULT NULL,
  `skill` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `buildskills_FK` (`build_id`),
  CONSTRAINT `buildskills_FK` FOREIGN KEY (`build_id`) REFERENCES `builds` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=49 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `buildskills`
--

LOCK TABLES `buildskills` WRITE;
/*!40000 ALTER TABLE `buildskills` DISABLE KEYS */;
/*!40000 ALTER TABLE `buildskills` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `buildstats`
--

DROP TABLE IF EXISTS `buildstats`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `buildstats` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `build_id` int(10) unsigned DEFAULT NULL,
  `stat` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `buildstats_FK` (`build_id`),
  CONSTRAINT `buildstats_FK` FOREIGN KEY (`build_id`) REFERENCES `builds` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=48 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `buildstats`
--

LOCK TABLES `buildstats` WRITE;
/*!40000 ALTER TABLE `buildstats` DISABLE KEYS */;
/*!40000 ALTER TABLE `buildstats` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `client`
--

DROP TABLE IF EXISTS `client`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `client` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(255) DEFAULT NULL,
  `password` varchar(100) DEFAULT NULL,
  `created_at` date DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_key` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1 COLLATE=latin1_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `client`
--

LOCK TABLES `client` WRITE;
/*!40000 ALTER TABLE `client` DISABLE KEYS */;
/*!40000 ALTER TABLE `client` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `client_session`
--

DROP TABLE IF EXISTS `client_session`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `client_session` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `session_token` varchar(255) DEFAULT NULL,
  `client_id` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `clientsessionUN` (`client_id`),
  CONSTRAINT `client_session_FK` FOREIGN KEY (`client_id`) REFERENCES `client` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=latin1 COLLATE=latin1_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `client_session`
--

LOCK TABLES `client_session` WRITE;
/*!40000 ALTER TABLE `client_session` DISABLE KEYS */;
/*!40000 ALTER TABLE `client_session` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `skill_values`
--

DROP TABLE IF EXISTS `skill_values`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `skill_values` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `skill_id` int(10) unsigned DEFAULT NULL,
  `skill_value` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `skill_values_FK` (`skill_id`),
  CONSTRAINT `skill_values_FK` FOREIGN KEY (`skill_id`) REFERENCES `buildskills` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=47 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `skill_values`
--

LOCK TABLES `skill_values` WRITE;
/*!40000 ALTER TABLE `skill_values` DISABLE KEYS */;
/*!40000 ALTER TABLE `skill_values` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `stat_values`
--

DROP TABLE IF EXISTS `stat_values`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `stat_values` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `stat_id` int(10) unsigned DEFAULT NULL,
  `stat_value` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `stat_values_FK` (`stat_id`),
  CONSTRAINT `stat_values_FK` FOREIGN KEY (`stat_id`) REFERENCES `buildstats` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `stat_values`
--

LOCK TABLES `stat_values` WRITE;
/*!40000 ALTER TABLE `stat_values` DISABLE KEYS */;
/*!40000 ALTER TABLE `stat_values` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'UDR'
--
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'IGNORE_SPACE,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
/*!50003 DROP PROCEDURE IF EXISTS `client_login` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_general_ci */ ;
DELIMITER ;;
CREATE DEFINER=`Cameron`@`localhost` PROCEDURE `client_login`(usernameInput varchar(255), passwordInput varchar(255), tokenInput varchar(500))
BEGIN
	declare clientId int unsigned;
    set clientId = (select id from client where username=usernameInput);
    if clientId is not null and passwordInput is not null then
      delete from client_session where client_id = clientId;
      if tokenInput is not null then 
	     insert into client_session (client_id,session_token) values (clientId,tokenInput);
	     select client_id, convert(session_token using "utf8") as session_token from client_session where client_id = clientId and session_token = tokenInput;
	     commit;
	  end if;
	end if;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'IGNORE_SPACE,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
/*!50003 DROP PROCEDURE IF EXISTS `client_signup` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_general_ci */ ;
DELIMITER ;;
CREATE DEFINER=`Cameron`@`localhost` PROCEDURE `client_signup`(usernameInput varchar(255), passwordInput varchar(255))
    MODIFIES SQL DATA
BEGIN
	if usernameInput is not null and passwordInput is not null then
	insert into client (username, password, created_at)
	values (usernameInput, passwordInput, now());
    commit;
	end if;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'IGNORE_SPACE,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
/*!50003 DROP PROCEDURE IF EXISTS `get_hpw` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_unicode_ci */ ;
DELIMITER ;;
CREATE DEFINER=`Cameron`@`localhost` PROCEDURE `get_hpw`(usernameInput varchar(255))
BEGIN
	select convert(password using "utf8") as password from client where username=usernameInput;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'IGNORE_SPACE,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
/*!50003 DROP PROCEDURE IF EXISTS `insert_build` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_general_ci */ ;
DELIMITER ;;
CREATE DEFINER=`Cameron`@`localhost` PROCEDURE `insert_build`(title_input varchar(255), client_id_input int unsigned, token_input varchar(500))
    MODIFIES SQL DATA
BEGIN
	
	declare client_id_checker varchar(255);
	
	if title_input is not null and client_id_input is not null then
	    if token_input is not null then
	       select client_id into client_id_checker
	       from client_session
	       where session_token = token_input;
	       if client_id_checker = client_id_input then
	          insert into builds (client_id, title) 
	          values (client_id_input, title_input); 
	          select LAST_INSERT_ID();
	          commit;
	       end if;
	    end if;
    end if;	   
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'IGNORE_SPACE,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
/*!50003 DROP PROCEDURE IF EXISTS `save_feats` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_general_ci */ ;
DELIMITER ;;
CREATE DEFINER=`Cameron`@`localhost` PROCEDURE `save_feats`(build_id_input int unsigned, token_input varchar(255), client_id_input int unsigned, feat_input varchar (100))
    MODIFIES SQL DATA
BEGIN
		declare token_checker varchar(255);
	    declare build_checker int unsigned;
	   
	   if token_input is not null and client_id_input is not null then
	           select session_token into token_checker
	           from client_session
	           where client_id = client_id_input;
	           if token_checker = token_input then
	                   if build_id_input is not null then
	                      select id into build_checker
	                 	  from builds 
                    	  where id = build_id_input;
                    	  if build_checker = build_id_input and feat_input is not null then
                    	      INSERT INTO buildfeats (build_id, feat)
                    	      values (build_id_input, feat_input);
                    	      commit;
                    	  end if;
                    end if;
               end if;
       end if;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'IGNORE_SPACE,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
/*!50003 DROP PROCEDURE IF EXISTS `save_skills` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_general_ci */ ;
DELIMITER ;;
CREATE DEFINER=`Cameron`@`localhost` PROCEDURE `save_skills`(build_id_input int unsigned, skill_input varchar(100), skill_value_input varchar(255),token_input varchar(255), client_id_input int unsigned)
BEGIN
	declare token_checker varchar(255);
	declare build_checker int unsigned;
    declare current_skill_id int unsigned;
	if build_id_input is not null THEN
	   select id into build_checker
	   from builds 
	   where id = build_id_input;
	   if build_checker is not null then
	      if token_input is not null and client_id_input is not null THEN
	         select session_token into token_checker 
	         from client_session
	         where client_id = client_id_input;
	         if token_checker = token_input and skill_input is not NULL then
	           insert into buildskills (build_id, skill)
	           values (build_id_input, skill_input);
	           commit;
	           select LAST_INSERT_ID() into current_skill_id
	           from buildskills
	           where skill = skill_input and build_id = build_id_input;
	          if current_skill_id is not null then
	              insert into skill_values (skill_id, skill_value)
	              values (current_skill_id, skill_value_input);
	              commit;
	             end if;   
	         end if;
	      end if;
	  end if;
	end if;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'IGNORE_SPACE,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
/*!50003 DROP PROCEDURE IF EXISTS `save_stats` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_general_ci */ ;
DELIMITER ;;
CREATE DEFINER=`Cameron`@`localhost` PROCEDURE `save_stats`(build_id_input int unsigned, stat_input varchar(100), stat_value_input varchar(255),token_input varchar(255), client_id_input int unsigned)
    MODIFIES SQL DATA
BEGIN
	
	declare token_checker varchar(255);
	declare build_checker int unsigned;
    declare current_stat_id int unsigned;
	if build_id_input is not null THEN
	   select id into build_checker
	   from builds 
	   where id = build_id_input;
	   if build_checker = build_id_input then
	      if token_input is not null and client_id_input is not null THEN
	         select session_token into token_checker 
	         from client_session
	         where client_id = client_id_input;
	         if token_checker = token_input and stat_input is not NULL then
	           insert into buildstats (build_id, stat)
	           values (build_id_input, stat_input);
	           commit;
	           select LAST_INSERT_ID() into current_stat_id
	           from buildstats
	           where stat = stat_input and build_id = build_id_input;
	          if current_stat_id is not null then
	              insert into stat_values (stat_id, stat_value)
	              values (current_stat_id, stat_value_input);
	              commit;
	             end if;   
	         end if;
	      end if;
	  end if;
	end if;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-11-22  6:20:23
